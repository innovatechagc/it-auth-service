package models

// Auth models - Modelos relacionados con autenticación
type LoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type LoginResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}

type LogoutRequest struct {
	SessionID string `json:"session_id" validate:"required"`
	IDToken   string `json:"id_token,omitempty"`
}

type AuthStatusRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type AuthStatusResponse struct {
	Authenticated bool  `json:"authenticated"`
	User          *User `json:"user,omitempty"`
	ExpiresAt     int64 `json:"expires_at,omitempty"`
}

// User model básico para auth service
type User struct {
	FirebaseID    string `json:"firebase_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Username      string `json:"username"`
	Status        string `json:"status"`
}