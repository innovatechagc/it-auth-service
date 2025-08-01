package models

import "time"

// RevokedToken representa un token JWT revocado
type RevokedToken struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenHash string    `json:"-" gorm:"uniqueIndex;not null"` // Hash del token para seguridad
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Reason    string    `json:"reason" gorm:"default:logout"` // logout, security, expired, etc.
	RevokedAt time.Time `json:"revoked_at" gorm:"autoCreateTime"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"` // Para limpiar tokens expirados
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	
	// Relación con User
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserSession representa una sesión de usuario para auditoría
type UserSession struct {
	ID         string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string     `json:"user_id" gorm:"not null;index"`
	TokenHash  string     `json:"-" gorm:"not null"` // Hash del JWT
	LoginAt    time.Time  `json:"login_at" gorm:"autoCreateTime"`
	LogoutAt   *time.Time `json:"logout_at,omitempty"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	Provider   string     `json:"provider"` // google.com, facebook.com, etc.
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	LastSeenAt time.Time  `json:"last_seen_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TokenValidationRequest para validar tokens
type TokenValidationRequest struct {
	Token string `json:"token" validate:"required"`
}

// TokenValidationResponse respuesta de validación
type TokenValidationResponse struct {
	Valid     bool   `json:"valid"`
	Reason    string `json:"reason,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}