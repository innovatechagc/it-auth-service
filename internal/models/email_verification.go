package models

import "time"

// EmailVerification representa el estado de verificación de email
type EmailVerification struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint       `json:"user_id" gorm:"not null;index"`
	FirebaseID       string     `json:"firebase_id" gorm:"size:128;not null;index"`
	Email            string     `json:"email" gorm:"size:255;not null;index"`
	IsVerified       bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	VerificationCode string     `json:"-" gorm:"size:6"` // No exponer en JSON
	CodeExpiresAt    *time.Time `json:"-"`               // No exponer en JSON
	AttemptsCount    int        `json:"attempts_count" gorm:"default:0"`
	LastAttemptAt    *time.Time `json:"last_attempt_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User (comentada porque User no tiene ID como primary key)
	// User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SendVerificationEmailRequest representa una solicitud para enviar email de verificación
type SendVerificationEmailRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Language string `json:"language,omitempty" validate:"omitempty,oneof=es en fr de it pt"`
}

// VerifyEmailRequest representa una solicitud para verificar email
type VerifyEmailRequest struct {
	IDToken          string `json:"id_token" validate:"required"`
	VerificationCode string `json:"verification_code,omitempty" validate:"omitempty,len=6,numeric"`
}

// VerifyEmailResponse representa la respuesta de verificación de email
type VerifyEmailResponse struct {
	Success       bool                `json:"success"`
	User          *User              `json:"user,omitempty"`
	Verification  *EmailVerification `json:"verification,omitempty"`
	Message       string             `json:"message"`
}

// EmailVerificationStatusRequest representa una solicitud para verificar el estado del email
type EmailVerificationStatusRequest struct {
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	FirebaseID string `json:"firebase_id,omitempty" validate:"omitempty,min=1"`
	UserID     int    `json:"user_id,omitempty" validate:"omitempty,min=1"`
}

// EmailVerificationStatusResponse representa la respuesta del estado de verificación
type EmailVerificationStatusResponse struct {
	EmailVerified    bool                `json:"email_verified"`
	Email           string              `json:"email"`
	UserID          int                 `json:"user_id,omitempty"`
	FirebaseID      string              `json:"firebase_id,omitempty"`
	Verification    *EmailVerification  `json:"verification,omitempty"`
	CanResend       bool                `json:"can_resend"`
	NextResendAt    *time.Time          `json:"next_resend_at,omitempty"`
	AttemptsLeft    int                 `json:"attempts_left"`
}

// ResendVerificationEmailRequest representa una solicitud para reenviar email de verificación
type ResendVerificationEmailRequest struct {
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	FirebaseID string `json:"firebase_id,omitempty" validate:"omitempty,min=1"`
	IDToken    string `json:"id_token,omitempty" validate:"omitempty,min=1"`
}

// VerifyEmailWithCodeRequest representa una solicitud para verificar email con código
type VerifyEmailWithCodeRequest struct {
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" validate:"required,len=6,numeric"`
}

// EmailVerificationSettings representa la configuración de verificación de email
type EmailVerificationSettings struct {
	MaxAttempts           int           `json:"max_attempts"`
	CodeExpirationTime    time.Duration `json:"code_expiration_time"`
	ResendCooldownTime    time.Duration `json:"resend_cooldown_time"`
	RequireVerification   bool          `json:"require_verification"`
	AutoVerifyDomains     []string      `json:"auto_verify_domains"`
	BlockedDomains        []string      `json:"blocked_domains"`
}

// EmailTemplate representa una plantilla de email
type EmailTemplate struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Language    string                 `json:"language"`
	Subject     string                 `json:"subject"`
	HTMLContent string                 `json:"html_content"`
	TextContent string                 `json:"text_content"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}