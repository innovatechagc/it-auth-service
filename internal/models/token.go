package models

// Token models - Modelos relacionados con tokens
type TokenVerifyRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type TokenVerifyResponse struct {
	Valid     bool  `json:"valid"`
	User      *User `json:"user,omitempty"`
	ExpiresAt int64 `json:"expires_at"`
	IssuedAt  int64 `json:"issued_at"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenRefreshResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type TokenRevokeRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type CustomTokenRequest struct {
	UID    string                 `json:"uid" validate:"required"`
	Claims map[string]interface{} `json:"claims,omitempty"`
}

type CustomTokenResponse struct {
	CustomToken string `json:"custom_token"`
	ExpiresIn   int64  `json:"expires_in"`
}