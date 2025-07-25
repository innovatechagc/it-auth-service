package middleware

import (
	"context"
	"net/http"
	"strings"

	"it-auth-service/internal/logger"
	"it-auth-service/pkg/firebase"
)

type AuthMiddleware struct {
	firebaseAuth *firebase.Auth
}

func NewAuthMiddleware(firebaseAuth *firebase.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuth: firebaseAuth,
	}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.GetLogger().Warn("Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.GetLogger().Warn("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Verificar token con Firebase
		decodedToken, err := a.firebaseAuth.VerifyIDToken(context.Background(), token)
		if err != nil {
			logger.GetLogger().WithError(err).Warn("Invalid Firebase token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Agregar información del usuario al contexto
		ctx := context.WithValue(r.Context(), "user_id", decodedToken.UID)
		ctx = context.WithValue(ctx, "user_email", decodedToken.Claims["email"])
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware opcional para endpoints públicos
func (a *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if decodedToken, err := a.firebaseAuth.VerifyIDToken(context.Background(), token); err == nil {
					ctx := context.WithValue(r.Context(), "user_id", decodedToken.UID)
					ctx = context.WithValue(ctx, "user_email", decodedToken.Claims["email"])
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}