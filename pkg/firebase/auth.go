package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Auth struct {
	client    *auth.Client
	projectID string
}

func NewAuth(serviceAccountPath string) (*Auth, error) {
	// Verificar si el archivo existe
	if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("firebase service account file not found: %s", serviceAccountPath)
	}

	// Configurar Firebase
	opt := option.WithCredentialsFile(serviceAccountPath)
	config := &firebase.Config{
		ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
	}
	
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	log.Printf("Firebase Auth initialized successfully for project: %s", config.ProjectID)
	return &Auth{
		client:    client,
		projectID: config.ProjectID,
	}, nil
}

func (a *Auth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}
	return token, nil
}

func (a *Auth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := a.client.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (a *Auth) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	user, err := a.client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (a *Auth) RevokeRefreshTokens(ctx context.Context, uid string) error {
	err := a.client.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}
	return nil
}

func (a *Auth) CreateCustomToken(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	token, err := a.client.CustomToken(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token: %w", err)
	}
	return token, nil
}

func (a *Auth) CreateCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	token, err := a.client.CustomTokenWithClaims(ctx, uid, claims)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token with claims: %w", err)
	}
	return token, nil
}

func (a *Auth) UpdateUser(ctx context.Context, uid string, user *auth.UserToUpdate) (*auth.UserRecord, error) {
	updatedUser, err := a.client.UpdateUser(ctx, uid, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (a *Auth) DeleteUser(ctx context.Context, uid string) error {
	err := a.client.DeleteUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (a *Auth) GetProjectID() string {
	return a.projectID
}

// GetAuthClient obtiene un cliente de Firebase Auth usando las credenciales del entorno
func GetAuthClient() (*auth.Client, error) {
	ctx := context.Background()
	
	// Intentar usar las credenciales por defecto de Google Cloud
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
	})
	
	if err != nil {
		// Si falla, intentar con el archivo de credenciales
		serviceAccountPath := "firebase-service-account.json"
		if path := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH"); path != "" {
			serviceAccountPath = path
		}
		
		if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("firebase credentials not found. Set GOOGLE_APPLICATION_CREDENTIALS or provide firebase-service-account.json")
		}
		
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, &firebase.Config{
			ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
		}, opt)
		
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
		}
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	return client, nil
}