package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"it-auth-service/internal/config"
	"it-auth-service/internal/database"
	"it-auth-service/internal/handlers"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/services"
)

type Server struct {
	config              *config.Config
	router              *gin.Engine
	firebaseAuthService *services.FirebaseAuthService
	userService         *services.UserService
	tokenService        *services.TokenService
}

func NewServer(cfg *config.Config) (*Server, error) {
	log := logger.GetLogger()

	// Configurar Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Conectar a la base de datos
	if err := database.Connect(*cfg); err != nil {
		log.WithError(err).Error("Database connection failed")
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Ejecutar migraciones
	if err := database.AutoMigrate(); err != nil {
		log.WithError(err).Error("Database migration failed")
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	// Obtener conexión a la base de datos
	db := database.GetDB()

	// Inicializar servicios
	userService := services.NewUserService(db)
	tokenService := services.NewTokenService(db)
	firebaseAuthService, err := services.NewFirebaseAuthService(cfg, userService, tokenService)
	if err != nil {
		log.WithError(err).Error("Firebase Auth service initialization failed")
		return nil, fmt.Errorf("firebase auth service initialization failed: %w", err)
	}

	// Crear router de Gin
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar CORS
	var corsConfig cors.Config

	if cfg.Environment == "development" {
		// Para desarrollo, permitir todos los orígenes
		corsConfig = cors.Config{
			AllowAllOrigins: true,
			AllowMethods: []string{
				"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
			},
			AllowHeaders: []string{
				"Origin", "Content-Type", "Accept", "Authorization", 
				"X-Requested-With", "X-User-ID", "X-Is-Admin",
			},
			ExposeHeaders: []string{
				"Content-Length", "Content-Type",
			},
			MaxAge: 12 * time.Hour,
		}
	} else {
		// Para producción, orígenes específicos
		corsConfig = cors.Config{
			AllowOrigins: []string{
				"https://your-frontend-domain.com",
			},
			AllowMethods: []string{
				"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
			},
			AllowHeaders: []string{
				"Origin", "Content-Type", "Accept", "Authorization", 
				"X-Requested-With", "X-User-ID", "X-Is-Admin",
			},
			ExposeHeaders: []string{
				"Content-Length", "Content-Type",
			},
			AllowCredentials: true,
			MaxAge:          12 * time.Hour,
		}
	}

	router.Use(cors.New(corsConfig))

	server := &Server{
		config:              cfg,
		router:              router,
		firebaseAuthService: firebaseAuthService,
		userService:         userService,
		tokenService:        tokenService,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// Configurar las rutas usando nuestros handlers de Gin
	handlers.SetupRoutes(s.router, s.firebaseAuthService, s.userService, s.tokenService)
}

func (s *Server) Start() error {
	log := logger.GetLogger()
	
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.WithField("address", addr).Info("Starting Auth Service server")
	
	return s.router.Run(addr)
}

func (s *Server) Close() error {
	return database.Close()
}