package middleware

import (
	"net/http"
	"strings"

	"it-auth-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTAuthMiddleware struct {
	jwtSecret string
	logger    *logrus.Logger
}

func NewJWTAuthMiddleware(cfg *config.Config, logger *logrus.Logger) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		jwtSecret: cfg.JWTSecret,
		logger:    logger,
	}
}

// RequireAuth middleware que requiere autenticación JWT
func (j *JWTAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			j.logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Verificar formato Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			j.logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verificar y parsear el token JWT
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			j.logger.WithError(err).Warn("Invalid JWT token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		// Extraer claims
		if claims, ok := token.Claims.(*JWTClaims); ok {
			// Agregar información del usuario al contexto
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Next()
		} else {
			j.logger.Error("Failed to extract JWT claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token claims",
			})
			c.Abort()
			return
		}
	}
}

// OptionalAuth middleware opcional para endpoints públicos
func (j *JWTAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(j.jwtSecret), nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(*JWTClaims); ok {
						c.Set("user_id", claims.UserID)
						c.Set("user_email", claims.Email)
					}
				}
			}
		}
		c.Next()
	}
}
