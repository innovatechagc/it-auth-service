package repositories

import (
	"it-auth-service/internal/models"
	"time"
)

// AuthRepositoryInterface define los métodos para autenticación
type AuthRepositoryInterface interface {
	// Session management
	CreateSession(userID, sessionID string, expiresAt time.Time) error
	GetSession(sessionID string) (*AuthSession, error)
	DeleteSession(sessionID string) error
	DeleteUserSessions(userID string) error
	
	// Login tracking
	RecordLoginAttempt(email, ip, userAgent string, success bool, failureReason string) error
	GetLoginAttempts(email string, limit int) ([]*LoginAttempt, error)
}

// EmailVerificationRepositoryInterface define los métodos para verificación de email
type EmailVerificationRepositoryInterface interface {
	Create(verification *models.EmailVerification) error
	GetByEmail(email string) (*models.EmailVerification, error)
	GetByFirebaseID(firebaseID string) (*models.EmailVerification, error)
	Update(verification *models.EmailVerification) error
	Delete(id uint) error
	
	// Verification operations
	MarkAsVerified(email string) error
	IncrementAttempts(email string) error
	CanResendCode(email string) (bool, error)
}

// PasswordResetRepositoryInterface define los métodos para reset de contraseña
type PasswordResetRepositoryInterface interface {
	Create(token *models.PasswordResetToken) error
	GetByToken(token string) (*models.PasswordResetToken, error)
	GetByCode(code string) (*models.PasswordResetToken, error)
	GetByEmail(email string) (*models.PasswordResetToken, error)
	Update(token *models.PasswordResetToken) error
	Delete(id uint) error
	
	// Token operations
	MarkAsUsed(token string) error
	CleanupExpiredTokens() error
}

// Helper structs
type AuthSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginAttempt struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Success       bool      `json:"success"`
	FailureReason string    `json:"failure_reason,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}