package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"it-auth-service/internal/logger"
	"it-auth-service/internal/models"
)

type UserService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:     db,
		logger: logger.GetLogger(),
	}
}

// GetUserByFirebaseID busca un usuario por su Firebase ID
func (s *UserService) GetUserByFirebaseID(ctx context.Context, firebaseID string) (*models.User, error) {
	var user models.User
	
	err := s.db.WithContext(ctx).Where("firebase_id = ?", firebaseID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		s.logger.WithError(err).Error("Failed to get user by Firebase ID")
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// GetUserByEmail busca un usuario por su email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		s.logger.WithError(err).Error("Failed to get user by email")
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// GetUserByID busca un usuario por su ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		s.logger.WithError(err).Error("Failed to get user by ID")
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Verificar que no exista un usuario con el mismo Firebase ID
	existingUser, err := s.GetUserByFirebaseID(ctx, user.FirebaseID)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with Firebase ID already exists")
	}

	// Verificar que no exista un usuario con el mismo email
	existingUser, err = s.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	// Crear el usuario
	err = s.db.WithContext(ctx).Create(user).Error
	if err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":     user.ID,
		"firebase_id": user.FirebaseID,
		"email":       user.Email,
	}).Info("User created successfully")

	return user, nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	err := s.db.WithContext(ctx).Save(user).Error
	if err != nil {
		s.logger.WithError(err).Error("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":     user.ID,
		"firebase_id": user.FirebaseID,
	}).Info("User updated successfully")

	return nil
}

// DeleteUser elimina un usuario (soft delete)
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("status", "deleted").Error
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// ListUsers lista usuarios con paginación
func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Contar total de usuarios activos
	err := s.db.WithContext(ctx).Model(&models.User{}).Where("status != ?", "deleted").Count(&total).Error
	if err != nil {
		s.logger.WithError(err).Error("Failed to count users")
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Obtener usuarios con paginación
	offset := (page - 1) * limit
	err = s.db.WithContext(ctx).
		Where("status != ?", "deleted").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		s.logger.WithError(err).Error("Failed to list users")
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// GetUserProfile obtiene el perfil completo de un usuario
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Status == "deleted" {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// UpdateUserProfile actualiza el perfil de un usuario
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, updates map[string]interface{}) (*models.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Status == "deleted" {
		return nil, fmt.Errorf("user not found")
	}

	// Aplicar actualizaciones permitidas
	allowedFields := map[string]bool{
		"username":   true,
		"first_name": true,
		"last_name":  true,
		"photo_url":  true,
	}

	updateData := make(map[string]interface{})
	for field, value := range updates {
		if allowedFields[field] {
			updateData[field] = value
		}
	}

	if len(updateData) == 0 {
		return user, nil
	}

	err = s.db.WithContext(ctx).Model(user).Updates(updateData).Error
	if err != nil {
		s.logger.WithError(err).Error("Failed to update user profile")
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Recargar el usuario actualizado
	return s.GetUserByID(ctx, userID)
}