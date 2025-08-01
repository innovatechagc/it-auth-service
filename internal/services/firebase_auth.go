package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"it-auth-service/internal/config"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
	"it-auth-service/pkg/firebase"
)

type FirebaseAuthService struct {
	firebaseClient *auth.Client
	config         *config.Config
	userService    *UserService
	tokenService   *TokenService
	logger         *logrus.Logger
}

func NewFirebaseAuthService(cfg *config.Config, userService *UserService, tokenService *TokenService) (*FirebaseAuthService, error) {
	firebaseClient, err := firebase.GetAuthClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase client: %w", err)
	}

	return &FirebaseAuthService{
		firebaseClient: firebaseClient,
		config:         cfg,
		userService:    userService,
		tokenService:   tokenService,
		logger:         logger.GetLogger(),
	}, nil
}

// VerifyFirebaseToken verifica un token de Firebase y extrae la información del usuario
func (s *FirebaseAuthService) VerifyFirebaseToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := s.firebaseClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		s.logger.WithError(err).Error("Failed to verify Firebase token")
		return nil, fmt.Errorf("invalid Firebase token: %w", err)
	}

	return token, nil
}

// FirebaseLogin maneja el login con token de Firebase
func (s *FirebaseAuthService) FirebaseLogin(ctx context.Context, req *models.FirebaseLoginRequest) (*models.AuthResponseData, error) {
	// Verificar el token de Firebase
	token, err := s.VerifyFirebaseToken(ctx, req.FirebaseToken)
	if err != nil {
		return nil, err
	}

	// Buscar usuario existente por Firebase ID
	user, err := s.userService.GetUserByFirebaseID(ctx, token.UID)
	isNewUser := false

	if err != nil {
		// Si el usuario no existe, intentar buscarlo por email primero
		email := getStringFromClaims(token.Claims, "email")
		if email != "" {
			existingUser, emailErr := s.userService.GetUserByEmail(ctx, email)
			if emailErr == nil && existingUser != nil {
				// Usuario existe con el mismo email, actualizar Firebase ID
				existingUser.FirebaseID = token.UID
				if updateErr := s.userService.UpdateUser(ctx, existingUser); updateErr != nil {
					s.logger.WithError(updateErr).Warn("Failed to update user Firebase ID")
				}
				user = existingUser
			} else {
				// Usuario no existe, crear uno nuevo (autoprovisionamiento)
				s.logger.WithField("firebase_id", token.UID).Info("User not found, creating new user")
				
				user, err = s.createUserFromFirebaseToken(ctx, token, req.Provider)
				if err != nil {
					// Si falla por duplicado, intentar obtener el usuario existente
					if existingUser, getErr := s.userService.GetUserByFirebaseID(ctx, token.UID); getErr == nil {
						user = existingUser
					} else {
						return nil, fmt.Errorf("failed to create user: %w", err)
					}
				} else {
					isNewUser = true
				}
			}
		} else {
			return nil, fmt.Errorf("no email found in Firebase token")
		}
	}

	// Actualizar información del usuario si es necesario
	if err := s.updateUserFromFirebaseToken(ctx, user, token, req.Provider); err != nil {
		s.logger.WithError(err).Warn("Failed to update user information")
	}

	// Generar JWT interno
	jwtToken, err := s.generateInternalJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Actualizar timestamp de último login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userService.UpdateUser(ctx, user); err != nil {
		s.logger.WithError(err).Warn("Failed to update user last login timestamp")
	}

	// Crear sesión de usuario para auditoría
	// Nota: En un handler real, obtendrías IP y User-Agent del contexto HTTP
	// Por ahora, usamos valores por defecto
	_, err = s.tokenService.CreateSession(ctx, user.ID, jwtToken, "unknown", "unknown", req.Provider)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to create user session")
		// No es un error crítico, continuamos
	}

	return &models.AuthResponseData{
		Token:     jwtToken,
		User:      user,
		IsNewUser: isNewUser,
	}, nil
}

// FirebaseRegister maneja el registro con token de Firebase
func (s *FirebaseAuthService) FirebaseRegister(ctx context.Context, req *models.FirebaseRegisterRequest) (*models.AuthResponseData, error) {
	// Verificar el token de Firebase
	token, err := s.VerifyFirebaseToken(ctx, req.FirebaseToken)
	if err != nil {
		return nil, err
	}

	// Verificar si el usuario ya existe
	existingUser, err := s.userService.GetUserByFirebaseID(ctx, token.UID)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Crear nuevo usuario con datos adicionales
	user, err := s.createUserFromFirebaseTokenWithData(ctx, token, req.Provider, req.RegistrationData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generar JWT interno
	jwtToken, err := s.generateInternalJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &models.AuthResponseData{
		Token:     jwtToken,
		User:      user,
		IsNewUser: true,
	}, nil
}

// RefreshToken maneja la renovación de tokens
func (s *FirebaseAuthService) RefreshToken(ctx context.Context, req *models.FirebaseRefreshRequest) (*models.AuthResponseData, error) {
	// Verificar el token de Firebase
	token, err := s.VerifyFirebaseToken(ctx, req.FirebaseToken)
	if err != nil {
		return nil, err
	}

	// Buscar usuario
	user, err := s.userService.GetUserByFirebaseID(ctx, token.UID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Actualizar información del usuario
	if err := s.updateUserFromFirebaseToken(ctx, user, token, user.Provider); err != nil {
		s.logger.WithError(err).Warn("Failed to update user information")
	}

	// Generar nuevo JWT interno
	jwtToken, err := s.generateInternalJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &models.AuthResponseData{
		Token:     jwtToken,
		User:      user,
		IsNewUser: false,
	}, nil
}

// createUserFromFirebaseToken crea un usuario desde un token de Firebase
func (s *FirebaseAuthService) createUserFromFirebaseToken(ctx context.Context, token *auth.Token, provider string) (*models.User, error) {
	user := &models.User{
		FirebaseID:    token.UID,
		Email:         getStringFromClaims(token.Claims, "email"),
		EmailVerified: getBoolFromClaims(token.Claims, "email_verified"),
		FirstName:     getStringFromClaims(token.Claims, "name"),
		Provider:      provider,
		PhotoURL:      getStringFromClaims(token.Claims, "picture"),
		Status:        "active",
	}

	// Generar username si no existe
	if user.Username == "" {
		user.Username = s.generateUsernameFromEmail(user.Email)
	}

	return s.userService.CreateUser(ctx, user)
}

// createUserFromFirebaseTokenWithData crea un usuario con datos adicionales
func (s *FirebaseAuthService) createUserFromFirebaseTokenWithData(ctx context.Context, token *auth.Token, provider string, registrationData map[string]interface{}) (*models.User, error) {
	user := &models.User{
		FirebaseID:    token.UID,
		Email:         getStringFromClaims(token.Claims, "email"),
		EmailVerified: getBoolFromClaims(token.Claims, "email_verified"),
		Provider:      provider,
		PhotoURL:      getStringFromClaims(token.Claims, "picture"),
		Status:        "active",
	}

	// Aplicar datos adicionales del registro
	if username, ok := registrationData["username"].(string); ok && username != "" {
		user.Username = username
	} else {
		user.Username = s.generateUsernameFromEmail(user.Email)
	}

	if firstName, ok := registrationData["first_name"].(string); ok {
		user.FirstName = firstName
	} else {
		user.FirstName = getStringFromClaims(token.Claims, "name")
	}

	if lastName, ok := registrationData["last_name"].(string); ok {
		user.LastName = lastName
	}

	return s.userService.CreateUser(ctx, user)
}

// updateUserFromFirebaseToken actualiza la información del usuario desde Firebase
func (s *FirebaseAuthService) updateUserFromFirebaseToken(ctx context.Context, user *models.User, token *auth.Token, provider string) error {
	updated := false

	// Actualizar email verificado
	emailVerified := getBoolFromClaims(token.Claims, "email_verified")
	if user.EmailVerified != emailVerified {
		user.EmailVerified = emailVerified
		updated = true
	}

	// Actualizar foto de perfil
	photoURL := getStringFromClaims(token.Claims, "picture")
	if user.PhotoURL != photoURL && photoURL != "" {
		user.PhotoURL = photoURL
		updated = true
	}

	// Actualizar provider si es diferente
	if user.Provider != provider {
		user.Provider = provider
		updated = true
	}

	if updated {
		return s.userService.UpdateUser(ctx, user)
	}

	return nil
}

// generateInternalJWT genera un JWT interno para el usuario
func (s *FirebaseAuthService) generateInternalJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     user.ID, // Ahora es string (UUID)
		"firebase_id": user.FirebaseID,
		"email":       user.Email,
		"username":    user.Username,
		"provider":    user.Provider,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// generateUsernameFromEmail genera un username desde un email
func (s *FirebaseAuthService) generateUsernameFromEmail(email string) string {
	if email == "" {
		return fmt.Sprintf("user_%d", time.Now().Unix())
	}

	// Extraer la parte antes del @
	at := 0
	for i, char := range email {
		if char == '@' {
			at = i
			break
		}
	}

	if at > 0 {
		return email[:at]
	}

	return fmt.Sprintf("user_%d", time.Now().Unix())
}

// Funciones auxiliares para extraer datos de claims
func getStringFromClaims(claims map[string]interface{}, key string) string {
	if value, ok := claims[key].(string); ok {
		return value
	}
	return ""
}

func getBoolFromClaims(claims map[string]interface{}, key string) bool {
	if value, ok := claims[key].(bool); ok {
		return value
	}
	return false
}

// GetJWTSecret devuelve el secreto JWT para validación de tokens
func (s *FirebaseAuthService) GetJWTSecret() string {
	return s.config.JWTSecret
}