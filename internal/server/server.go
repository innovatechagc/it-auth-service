package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"it-auth-service/internal/config"
	"it-auth-service/internal/database"
	"it-auth-service/internal/handlers"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/middleware"
	"it-auth-service/pkg/firebase"
)

type Server struct {
	config       config.Config
	router       *mux.Router
	firebaseAuth *firebase.Auth
}

func NewServer(cfg config.Config) (*Server, error) {
	log := logger.GetLogger()

	// Conectar a la base de datos
	if err := database.Connect(cfg); err != nil {
		log.WithError(err).Warn("Database connection failed, some features will be disabled")
		// No es un error fatal para el auth service
	} else {
		// Ejecutar migraciones si la conexión fue exitosa
		if err := database.AutoMigrate(); err != nil {
			log.WithError(err).Warn("Database migration failed")
		}
	}

	// Inicializar Firebase Auth
	firebaseAuth, err := firebase.NewAuth("firebase-service-account.json")
	if err != nil {
		log.WithError(err).Warn("Firebase Auth initialization failed, some features will be disabled")
		// No es un error fatal, el servicio puede funcionar sin Firebase para algunos endpoints
	}

	server := &Server{
		config:       cfg,
		router:       mux.NewRouter(),
		firebaseAuth: firebaseAuth,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// Crear handlers
	authHandler := handlers.NewAuthHandler(s.firebaseAuth)
	tokenHandler := handlers.NewTokenHandler(s.firebaseAuth)

	// Crear middleware
	authMiddleware := middleware.NewAuthMiddleware(s.firebaseAuth)

	// Health check
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "auth-service"}`))
	}).Methods("GET")

	// Rutas públicas de autenticación
	s.router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	s.router.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	s.router.HandleFunc("/auth/status", authHandler.CheckAuthStatus).Methods("GET")
	s.router.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")

	// Rutas de tokens
	s.router.HandleFunc("/tokens/verify", tokenHandler.VerifyToken).Methods("POST")
	s.router.HandleFunc("/tokens/validate", tokenHandler.ValidateToken).Methods("POST")
	s.router.HandleFunc("/tokens/refresh", tokenHandler.RefreshToken).Methods("POST")

	// Rutas protegidas (requieren autenticación)
	protected := s.router.PathPrefix("/auth").Subrouter()
	protected.Use(authMiddleware.RequireAuth)
	
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/revoke-tokens", authHandler.RevokeAllTokens).Methods("POST")
	
	// Rutas protegidas de tokens
	protectedTokens := s.router.PathPrefix("/tokens").Subrouter()
	protectedTokens.Use(authMiddleware.RequireAuth)
	
	protectedTokens.HandleFunc("/revoke", tokenHandler.RevokeToken).Methods("POST")
	protectedTokens.HandleFunc("/revoke-all", tokenHandler.RevokeAllTokens).Methods("POST")
	protectedTokens.HandleFunc("/custom", tokenHandler.CreateCustomToken).Methods("POST")
}

func (s *Server) Start() error {
	log := logger.GetLogger()
	
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.WithField("address", addr).Info("Starting Auth Service server")
	
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) Close() error {
	return database.Close()
}