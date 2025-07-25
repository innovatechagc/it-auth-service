package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"it-auth-service/internal/config"
	"it-auth-service/internal/models"
)

var DB *gorm.DB

// Connect establece la conexión con la base de datos
func Connect(cfg config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate ejecuta las migraciones automáticas para tablas de autenticación
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Migrar solo modelos específicos de autenticación
	// Estas tablas son adicionales a las que ya existen en it-app_user
	err := DB.AutoMigrate(
		&models.EmailVerification{},
		&models.PasswordResetToken{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate auth database tables: %w", err)
	}

	log.Println("Auth service database migration completed successfully")
	return nil
}

// Close cierra la conexión a la base de datos
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}