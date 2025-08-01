package models

import "time"

// Auth models - Modelos relacionados con autenticación
type LoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type LoginResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}



type AuthStatusRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type AuthStatusResponse struct {
	Authenticated bool  `json:"authenticated"`
	User          *User `json:"user,omitempty"`
	ExpiresAt     int64 `json:"expires_at,omitempty"`
}

// User model completo para auth service
type User struct {
	ID            string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FirebaseID    string `json:"firebase_id" gorm:"uniqueIndex;not null"`
	Email         string `json:"email" gorm:"uniqueIndex;not null"`
	Username      string `json:"username" gorm:"uniqueIndex"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Provider      string `json:"provider"` // google.com, facebook.com, password
	PhotoURL      string `json:"photo_url"`
	Status        string     `json:"status" gorm:"default:active"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	LastLogoutAt  *time.Time `json:"last_logout_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// Firebase Login Request
type FirebaseLoginRequest struct {
	FirebaseToken string `json:"firebase_token" validate:"required"`
	Provider      string `json:"provider" validate:"required,oneof=google.com facebook.com password"`
}

// Firebase Register Request
type FirebaseRegisterRequest struct {
	FirebaseToken    string                 `json:"firebase_token" validate:"required"`
	Provider         string                 `json:"provider" validate:"required,oneof=google.com facebook.com password"`
	RegistrationData map[string]interface{} `json:"registration_data"`
}

// Firebase Refresh Token Request
type FirebaseRefreshRequest struct {
	FirebaseToken string `json:"firebase_token" validate:"required"`
	Provider      string `json:"provider" validate:"required,oneof=refresh"`
}

// Standard API Response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Auth Response Data
type AuthResponseData struct {
	Token     string `json:"token"`
	User      *User  `json:"user"`
	IsNewUser bool   `json:"isNewUser"`
}

// Logout Request
type LogoutRequest struct {
	Token string `json:"token,omitempty"`
}

// Logout Response
type LogoutResponse struct {
	Message string `json:"message"`
}