package models

import "time"

// Password Reset models - Modelos relacionados con reset de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Code        string `json:"code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentToken string `json:"current_token" validate:"required"`
	NewPassword  string `json:"new_password" validate:"required,min=6"`
}

type PasswordResetToken struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint       `json:"user_id" gorm:"not null;index"`
	FirebaseID string     `json:"firebase_id" gorm:"size:128;not null;index"`
	Email      string     `json:"email" gorm:"size:255;not null;index"`
	Token      string     `json:"-" gorm:"size:255;uniqueIndex"` // No exponer en JSON
	Code       string     `json:"-" gorm:"size:6"`               // No exponer en JSON
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt     *time.Time `json:"used_at,omitempty"`
	IsUsed     bool       `json:"is_used" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type PasswordResetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type PasswordStrengthCheck struct {
	IsValid      bool     `json:"is_valid"`
	Score        int      `json:"score"` // 0-4
	Feedback     []string `json:"feedback,omitempty"`
	Requirements struct {
		MinLength    bool `json:"min_length"`
		HasUppercase bool `json:"has_uppercase"`
		HasLowercase bool `json:"has_lowercase"`
		HasNumbers   bool `json:"has_numbers"`
		HasSymbols   bool `json:"has_symbols"`
	} `json:"requirements"`
}

type PasswordPolicy struct {
	MinLength        int      `json:"min_length"`
	RequireUppercase bool     `json:"require_uppercase"`
	RequireLowercase bool     `json:"require_lowercase"`
	RequireNumbers   bool     `json:"require_numbers"`
	RequireSymbols   bool     `json:"require_symbols"`
	ForbiddenWords   []string `json:"forbidden_words,omitempty"`
	MaxAge           int      `json:"max_age_days,omitempty"` // días
}