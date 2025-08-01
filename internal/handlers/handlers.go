package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
	"it-auth-service/internal/services"
)

type Handler struct {
	firebaseAuthService *services.FirebaseAuthService
	userService         *services.UserService
	tokenService        *services.TokenService
	logger              *logrus.Logger
}

func NewHandler(firebaseAuthService *services.FirebaseAuthService, userService *services.UserService, tokenService *services.TokenService) *Handler {
	return &Handler{
		firebaseAuthService: firebaseAuthService,
		userService:         userService,
		tokenService:        tokenService,
		logger:              logger.GetLogger(),
	}
}

func SetupRoutes(router *gin.Engine, firebaseAuthService *services.FirebaseAuthService, userService *services.UserService, tokenService *services.TokenService) {
	h := NewHandler(firebaseAuthService, userService, tokenService)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", h.HealthCheck)
		api.GET("/ready", h.ReadinessCheck)
		api.GET("/info", h.APIInfo)

		// Firebase Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/firebase-login", h.FirebaseLogin)
			auth.POST("/firebase-register", h.FirebaseRegister)
			auth.POST("/refresh-token", h.RefreshToken)
			auth.POST("/logout", h.Logout)
		}

		// User Management
		users := api.Group("/users")
		{
			users.GET("/profile", h.GetUserProfile)
			users.PUT("/profile", h.UpdateUserProfile)
			users.GET("", h.ListUsers) // Admin only
		}
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Verifica el estado del servicio
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "it-auth-service",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}
	
	c.JSON(http.StatusOK, response)
}

// ReadinessCheck godoc
// @Summary Readiness check endpoint
// @Description Verifica si el servicio está listo para recibir tráfico
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *Handler) ReadinessCheck(c *gin.Context) {
	response := map[string]interface{}{
		"ready":     true,
		"service":   "it-auth-service",
		"timestamp": time.Now().UTC(),
	}
	
	c.JSON(http.StatusOK, response)
}

// Ejemplo de handler comentado para testing
/*
// GetExample godoc
// @Summary Get example data
// @Description Obtiene datos de ejemplo
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /example [get]
func (h *Handler) GetExample(c *gin.Context) {
	// Implementación de ejemplo
	c.JSON(http.StatusOK, gin.H{
		"message": "Example data",
		"data":    []string{"item1", "item2", "item3"},
	})
}

// CreateExample godoc
// @Summary Create example data
// @Description Crea datos de ejemplo
// @Tags example
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Example data"
// @Success 201 {object} map[string]interface{}
// @Router /example [post]
func (h *Handler) CreateExample(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Implementación de ejemplo
	c.JSON(http.StatusCreated, gin.H{
		"message": "Example created",
		"data":    request,
	})
}
*/

// APIInfo godoc
// @Summary API information endpoint
// @Description Obtiene información general de la API
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /info [get]
func (h *Handler) APIInfo(c *gin.Context) {
	response := map[string]interface{}{
		"name":        "IT Auth Service",
		"version":     "1.0.0",
		"description": "Servicio de autenticación con Firebase",
		"endpoints": map[string]interface{}{
			"health":           "/api/v1/health",
			"ready":            "/api/v1/ready",
			"firebase_login":   "/api/v1/auth/firebase-login",
			"firebase_register": "/api/v1/auth/firebase-register",
			"refresh_token":    "/api/v1/auth/refresh-token",
			"user_profile":     "/api/v1/users/profile",
		},
		"timestamp": time.Now().UTC(),
	}
	
	c.JSON(http.StatusOK, response)
}

// FirebaseLogin godoc
// @Summary Firebase login endpoint
// @Description Autenticación con token de Firebase
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.FirebaseLoginRequest true "Firebase login data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/firebase-login [post]
func (h *Handler) FirebaseLogin(c *gin.Context) {
	var req models.FirebaseLoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for Firebase login")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar provider
	if req.Provider != "google.com" && req.Provider != "facebook.com" && req.Provider != "password" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid provider. Must be one of: google.com, facebook.com, password",
		})
		return
	}

	authData, err := h.firebaseAuthService.FirebaseLogin(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Firebase login failed")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication failed: " + err.Error(),
		})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id":    authData.User.ID,
		"email":      authData.User.Email,
		"provider":   authData.User.Provider,
		"is_new_user": authData.IsNewUser,
	}).Info("Firebase login successful")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    authData,
	})
}

// FirebaseRegister godoc
// @Summary Firebase register endpoint
// @Description Registro con token de Firebase
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.FirebaseRegisterRequest true "Firebase register data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/firebase-register [post]
func (h *Handler) FirebaseRegister(c *gin.Context) {
	var req models.FirebaseRegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for Firebase register")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar provider
	if req.Provider != "google.com" && req.Provider != "facebook.com" && req.Provider != "password" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid provider. Must be one of: google.com, facebook.com, password",
		})
		return
	}

	authData, err := h.firebaseAuthService.FirebaseRegister(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Firebase register failed")
		
		statusCode := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			statusCode = http.StatusConflict
		}
		
		c.JSON(statusCode, models.APIResponse{
			Success: false,
			Error:   "Registration failed: " + err.Error(),
		})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id":  authData.User.ID,
		"email":    authData.User.Email,
		"provider": authData.User.Provider,
	}).Info("Firebase register successful")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    authData,
	})
}

// RefreshToken godoc
// @Summary Refresh token endpoint
// @Description Renovación de token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.FirebaseRefreshRequest true "Refresh token data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/refresh-token [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req models.FirebaseRefreshRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for refresh token")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar provider
	if req.Provider != "refresh" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid provider. Must be 'refresh'",
		})
		return
	}

	authData, err := h.firebaseAuthService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Token refresh failed")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Token refresh failed: " + err.Error(),
		})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id": authData.User.ID,
		"email":   authData.User.Email,
	}).Info("Token refresh successful")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    authData,
	})
}

// GetUserProfile godoc
// @Summary Get user profile endpoint
// @Description Obtiene el perfil del usuario autenticado
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/profile [get]
func (h *Handler) GetUserProfile(c *gin.Context) {
	// TODO: Implementar middleware de autenticación JWT
	// Por ahora, simulamos obtener el user_id del token
	userID := c.GetHeader("X-User-ID") // Temporal para testing
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	user, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user profile")
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user": user,
		},
	})
}

// UpdateUserProfile godoc
// @Summary Update user profile endpoint
// @Description Actualiza el perfil del usuario autenticado
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Profile update data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/profile [put]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	// TODO: Implementar middleware de autenticación JWT
	// Por ahora, simulamos obtener el user_id del token
	userID := c.GetHeader("X-User-ID") // Temporal para testing
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		h.logger.WithError(err).Error("Invalid request body for profile update")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, updates)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update user profile")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update profile: " + err.Error(),
		})
		return
	}

	h.logger.WithField("user_id", userID).Info("User profile updated successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Profile updated successfully",
			"user":    user,
		},
	})
}

// ListUsers godoc
// @Summary List users endpoint (Admin only)
// @Description Lista usuarios con paginación
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	// TODO: Implementar middleware de autenticación y autorización
	// Por ahora, simulamos la verificación de admin
	isAdmin := c.GetHeader("X-Is-Admin") == "true" // Temporal para testing
	if !isAdmin {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Admin access required",
		})
		return
	}

	// Obtener parámetros de paginación
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list users")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to list users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// Logout godoc
// @Summary Logout endpoint
// @Description Cierra la sesión del usuario e invalida el token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LogoutRequest true "Logout data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	
	// El logout puede ser llamado sin body (solo con headers)
	if err := c.ShouldBindJSON(&req); err != nil {
		// Si no hay body, está bien, solo logueamos que no se pudo parsear
		h.logger.WithError(err).Debug("No JSON body provided for logout, proceeding anyway")
	}

	// Obtener token del header Authorization si existe
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		req.Token = authHeader[7:]
	}

	// Validar que tenemos un token
	if req.Token == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Token is required for logout",
		})
		return
	}

	// Parsear el token para obtener el user_id
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.firebaseAuthService.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid token claims",
		})
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID in token",
		})
		return
	}

	// Obtener información adicional
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 1. Revocar el token (agregarlo a la blacklist)
	if err := h.tokenService.RevokeToken(c.Request.Context(), req.Token, userID, "logout", ipAddress, userAgent); err != nil {
		h.logger.WithError(err).Error("Failed to revoke token")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to revoke token",
		})
		return
	}

	// 2. Terminar la sesión del usuario
	if err := h.tokenService.EndSession(c.Request.Context(), req.Token, userID); err != nil {
		h.logger.WithError(err).Warn("Failed to end user session")
		// No es un error crítico, continuamos
	}

	// 3. Actualizar el timestamp de último logout del usuario
	now := time.Now()
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err == nil {
		user.LastLogoutAt = &now
		if err := h.userService.UpdateUser(c.Request.Context(), user); err != nil {
			h.logger.WithError(err).Warn("Failed to update user last logout timestamp")
			// No es un error crítico, continuamos
		}
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"ip":      ipAddress,
		"reason":  "logout",
	}).Info("User logout successful - token revoked")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.LogoutResponse{
			Message: "Logout successful - token has been revoked",
		},
	})
}