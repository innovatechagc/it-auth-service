package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
)

type TokenService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{
		db:     db,
		logger: logger.GetLogger(),
	}
}

// RevokeToken revoca un token JWT
func (s *TokenService) RevokeToken(ctx context.Context, tokenString, userID, reason, ipAddress, userAgent string) error {
	// Parsear el token para obtener la fecha de expiración (sin validar la firma)
	token, _ := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// No validamos la firma aquí, solo queremos extraer los claims
		return nil, nil
	})
	
	// Incluso si hay error de validación, podemos extraer los claims
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		if exp, exists := (*claims)["exp"]; exists {
			if expFloat, ok := exp.(float64); ok {
				expiresAt := time.Unix(int64(expFloat), 0)
				
				// Crear hash del token para seguridad
				tokenHash := s.hashToken(tokenString)
				
				revokedToken := &models.RevokedToken{
					TokenHash: tokenHash,
					UserID:    userID,
					Reason:    reason,
					ExpiresAt: expiresAt,
					IPAddress: ipAddress,
					UserAgent: userAgent,
				}
				
				// Guardar token revocado
				if err := s.db.WithContext(ctx).Create(revokedToken).Error; err != nil {
					s.logger.WithError(err).Error("Failed to revoke token")
					return fmt.Errorf("failed to revoke token: %w", err)
				}
				
				s.logger.WithFields(map[string]interface{}{
					"user_id": userID,
					"reason":  reason,
					"ip":      ipAddress,
				}).Info("Token revoked successfully")
				
				return nil
			}
		}
	}
	
	return fmt.Errorf("invalid token format or missing expiration")
}

// IsTokenRevoked verifica si un token está revocado
func (s *TokenService) IsTokenRevoked(ctx context.Context, tokenString string) (bool, error) {
	tokenHash := s.hashToken(tokenString)
	
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.RevokedToken{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count).Error
	
	if err != nil {
		s.logger.WithError(err).Error("Failed to check token revocation status")
		return false, fmt.Errorf("failed to check token status: %w", err)
	}
	
	return count > 0, nil
}

// CreateSession crea una nueva sesión de usuario
func (s *TokenService) CreateSession(ctx context.Context, userID, tokenString, ipAddress, userAgent, provider string) (*models.UserSession, error) {
	tokenHash := s.hashToken(tokenString)
	
	session := &models.UserSession{
		UserID:    userID,
		TokenHash: tokenHash,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Provider:  provider,
		IsActive:  true,
	}
	
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		s.logger.WithError(err).Error("Failed to create user session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	
	s.logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"session_id": session.ID,
		"provider":   provider,
		"ip":         ipAddress,
	}).Info("User session created")
	
	return session, nil
}

// EndSession termina una sesión de usuario
func (s *TokenService) EndSession(ctx context.Context, tokenString, userID string) error {
	tokenHash := s.hashToken(tokenString)
	now := time.Now()
	
	// Actualizar sesión como inactiva
	err := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("token_hash = ? AND user_id = ? AND is_active = ?", tokenHash, userID, true).
		Updates(map[string]interface{}{
			"logout_at": &now,
			"is_active": false,
		}).Error
	
	if err != nil {
		s.logger.WithError(err).Error("Failed to end user session")
		return fmt.Errorf("failed to end session: %w", err)
	}
	
	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
	}).Info("User session ended")
	
	return nil
}

// UpdateLastSeen actualiza la última actividad de la sesión
func (s *TokenService) UpdateLastSeen(ctx context.Context, tokenString, userID string) error {
	tokenHash := s.hashToken(tokenString)
	
	err := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("token_hash = ? AND user_id = ? AND is_active = ?", tokenHash, userID, true).
		Update("last_seen_at", time.Now()).Error
	
	if err != nil {
		s.logger.WithError(err).Debug("Failed to update last seen")
		// No es un error crítico, solo log como debug
	}
	
	return nil
}

// CleanupExpiredTokens limpia tokens expirados de la base de datos
func (s *TokenService) CleanupExpiredTokens(ctx context.Context) error {
	now := time.Now()
	
	// Limpiar tokens revocados expirados
	result := s.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&models.RevokedToken{})
	
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to cleanup expired revoked tokens")
		return fmt.Errorf("failed to cleanup expired tokens: %w", result.Error)
	}
	
	// Limpiar sesiones inactivas antiguas (más de 30 días)
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	result = s.db.WithContext(ctx).
		Where("is_active = ? AND (logout_at < ? OR last_seen_at < ?)", false, thirtyDaysAgo, thirtyDaysAgo).
		Delete(&models.UserSession{})
	
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to cleanup old sessions")
		return fmt.Errorf("failed to cleanup old sessions: %w", result.Error)
	}
	
	s.logger.WithFields(map[string]interface{}{
		"revoked_tokens_cleaned": result.RowsAffected,
	}).Info("Token cleanup completed")
	
	return nil
}

// GetUserActiveSessions obtiene las sesiones activas de un usuario
func (s *TokenService) GetUserActiveSessions(ctx context.Context, userID string) ([]*models.UserSession, error) {
	var sessions []*models.UserSession
	
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_seen_at DESC").
		Find(&sessions).Error
	
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user active sessions")
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	
	return sessions, nil
}

// hashToken crea un hash SHA256 del token para almacenamiento seguro
func (s *TokenService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}