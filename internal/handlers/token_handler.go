package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
	"it-auth-service/internal/validator"
	"it-auth-service/pkg/firebase"
)

type TokenHandler struct {
	firebaseAuth *firebase.Auth
}

func NewTokenHandler(firebaseAuth *firebase.Auth) *TokenHandler {
	return &TokenHandler{
		firebaseAuth: firebaseAuth,
	}
}

// VerifyToken maneja POST /auth/verify-token
func (h *TokenHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.TokenVerifyRequest

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
		log.WithError(err).Warn("Validation failed for token verify request")
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

	// Crear respuesta con información del token
	response := models.TokenVerifyResponse{
		Valid: true,
		User: &models.User{
			FirebaseID:    token.UID,
			Email:         userRecord.Email,
			EmailVerified: userRecord.EmailVerified,
			Username:      userRecord.DisplayName,
			Status:        "active",
		},
		ExpiresAt: token.Expires,
		IssuedAt:  token.IssuedAt,
	}

	log.WithField("firebase_id", token.UID).Info("Token verified successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken maneja POST /auth/refresh-token
func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.TokenRefreshRequest

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
		log.WithError(err).Warn("Validation failed for token refresh request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// En Firebase, el refresh de tokens se maneja del lado del cliente
	// Aquí registramos el evento y enviamos una respuesta genérica
	
	log.Info("Token refresh requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Token refresh should be handled on the client side",
		"info":    "Use Firebase SDK to refresh tokens automatically",
	})
}

// RevokeToken maneja POST /auth/revoke-token
func (h *TokenHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.TokenRevokeRequest

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
		log.WithError(err).Warn("Validation failed for token revoke request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Obtener el usuario del contexto (debe estar autenticado)
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// En Firebase, podemos revocar todos los tokens de un usuario
	err := h.firebaseAuth.RevokeRefreshTokens(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to revoke tokens")
		http.Error(w, "Failed to revoke tokens", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", userID).Info("Tokens revoked successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All tokens revoked successfully",
	})
}

// CreateCustomToken maneja POST /tokens/custom
func (h *TokenHandler) CreateCustomToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.CustomTokenRequest

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
		log.WithError(err).Warn("Validation failed for custom token request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	log.WithField("uid", req.UID).Info("Custom token creation requested (requires Firebase Admin SDK)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Custom token creation requires Firebase Admin SDK implementation",
	})
}

// RevokeAllTokens maneja POST /tokens/revoke-all
func (h *TokenHandler) RevokeAllTokens(w http.ResponseWriter, r *http.Request) {
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

	// Revocar todos los tokens del usuario
	err := h.firebaseAuth.RevokeRefreshTokens(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to revoke all tokens")
		http.Error(w, "Failed to revoke tokens", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", userID).Info("All tokens revoked successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All tokens revoked successfully",
	})
}

// ValidateToken maneja POST /tokens/validate
func (h *TokenHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.TokenVerifyRequest

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
		log.WithError(err).Warn("Validation failed for token validate request")
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
		log.WithError(err).Warn("Token validation failed")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   false,
			"message": "Invalid token",
		})
		return
	}

	log.WithField("firebase_id", token.UID).Info("Token validated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":      true,
		"user_id":    token.UID,
		"expires_at": token.Expires,
		"issued_at":  token.IssuedAt,
		"message":    "Token is valid",
	})
}