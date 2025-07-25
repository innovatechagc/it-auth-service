package main

import (
	"os"

	"github.com/joho/godotenv"
	"it-auth-service/internal/config"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/server"
)

func main() {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		// No es un error fatal si no existe el archivo .env
		// Las variables de entorno del sistema tendrán precedencia
	}

	// Inicializar logger
	logger.Init()
	log := logger.GetLogger()

	// Cargar configuración
	cfg := config.LoadConfig()
	
	log.WithField("environment", cfg.Environment).WithField("port", cfg.Port).Info("Starting Auth Service")

	// Crear servidor
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.WithError(err).Fatal("Error creating server")
		os.Exit(1)
	}

	// Cleanup al cerrar
	defer func() {
		if err := srv.Close(); err != nil {
			log.WithError(err).Error("Error closing server")
		}
	}()

	// Iniciar servidor
	log.WithField("port", cfg.Port).Info("Auth Server initialized successfully")
	if err := srv.Start(); err != nil {
		log.WithError(err).Fatal("Error starting server")
		os.Exit(1)
	}
}