package firebase

import (
	"encoding/json"
	"fmt"
	"os"
)

// FirebaseConfig representa la configuración de Firebase
type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// LoadFirebaseConfig carga la configuración de Firebase desde el archivo JSON
func LoadFirebaseConfig(path string) (*FirebaseConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open firebase config file: %w", err)
	}
	defer file.Close()

	var config FirebaseConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode firebase config: %w", err)
	}

	return &config, nil
}

// ValidateConfig valida que la configuración de Firebase sea correcta
func ValidateConfig(config *FirebaseConfig) error {
	if config.ProjectID == "" {
		return fmt.Errorf("firebase project_id is required")
	}
	if config.ClientEmail == "" {
		return fmt.Errorf("firebase client_email is required")
	}
	if config.PrivateKey == "" {
		return fmt.Errorf("firebase private_key is required")
	}
	return nil
}

// GetProjectIDFromEnv obtiene el Project ID desde las variables de entorno
func GetProjectIDFromEnv() string {
	return os.Getenv("FIREBASE_PROJECT_ID")
}