package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
	"it-auth-service/internal/validator"
	"it-auth-service/pkg/firebase"
)

type AuthHandler struct {
	firebaseAuth *firebase.Auth
}

func NewAuthHandler(firebaseAuth *firebase.Auth) *AuthHandler {
	return &AuthHandler{
		firebaseAuth: firebaseAuth,
	}
}

// Login maneja POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.LoginRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for login request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar token de Firebase
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Firebase token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Crear respuesta de login
	response := models.LoginResponse{
		User: &models.User{
			FirebaseID:    token.UID,
			Email:         userRecord.Email,
			EmailVerified: userRecord.EmailVerified,
			Username:      userRecord.DisplayName,
			Status:        "active",
		},
		Message: "Login successful",
	}

	log.WithField("firebase_id", token.UID).Info("User logged in successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout maneja POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// En Firebase, el logout se maneja del lado del cliente
	// Aquí podríamos registrar el evento de logout
	
	log.Info("User logout requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}

// CheckAuthStatus maneja GET /auth/status
func (h *AuthHandler) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener token del header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "No token provided",
		})
		return
	}

	// Verificar formato Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid token format",
		})
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	token := parts[1]
	decodedToken, err := h.firebaseAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid token",
		})
		return
	}

	// Obtener información del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), decodedToken.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": true,
			"user_id":       decodedToken.UID,
			"expires_at":    decodedToken.Expires,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"firebase_id":    decodedToken.UID,
			"email":          userRecord.Email,
			"email_verified": userRecord.EmailVerified,
			"display_name":   userRecord.DisplayName,
		},
		"expires_at": decodedToken.Expires,
	})
}

// RefreshToken maneja POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// En Firebase, el refresh de tokens se maneja del lado del cliente
	log.Info("Token refresh requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Token refresh should be handled on the client side using Firebase SDK",
	})
}

// GetProfile maneja GET /auth/profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto (debe estar autenticado)
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to get user profile from Firebase")
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	profile := map[string]interface{}{
		"firebase_id":    userRecord.UID,
		"email":          userRecord.Email,
		"email_verified": userRecord.EmailVerified,
		"display_name":   userRecord.DisplayName,
		"photo_url":      userRecord.PhotoURL,
		"phone_number":   userRecord.PhoneNumber,
		"disabled":       userRecord.Disabled,
		"created_at":     userRecord.UserMetadata.CreationTimestamp,
		"last_login_at":  userRecord.UserMetadata.LastLogInTimestamp,
	}

	log.WithField("user_id", userID).Info("User profile retrieved")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    profile,
		"message": "Profile retrieved successfully",
	})
}

// RevokeAllTokens maneja POST /auth/revoke-tokens
func (h *AuthHandler) RevokeAllTokens(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Revocar todos los tokens del usuario
	err := h.firebaseAuth.RevokeRefreshTokens(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to revoke tokens")
		http.Error(w, "Failed to revoke tokens", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", userID).Info("All tokens revoked successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All tokens revoked successfully",
	})
}