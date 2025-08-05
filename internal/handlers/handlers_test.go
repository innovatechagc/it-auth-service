package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Mock simple para testing
type MockServices struct {
	logger *logrus.Logger
}

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Crear logger para testing
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Crear handler simple para testing
	handler := &Handler{
		logger: logger,
	}

	// Agregar ruta de health check directamente
	router.GET("/health", handler.HealthCheck)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Assertions básicas
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}
}

func TestReadinessCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Crear logger para testing
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Crear handler simple para testing
	handler := &Handler{
		logger: logger,
	}

	// Agregar ruta de readiness check directamente
	router.GET("/ready", handler.ReadinessCheck)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ready", nil)
	router.ServeHTTP(w, req)

	// Assertions básicas
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}
}

func TestAPIInfo(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Crear logger para testing
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Crear handler simple para testing
	handler := &Handler{
		logger: logger,
	}

	// Agregar ruta de API info directamente
	router.GET("/info", handler.APIInfo)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	router.ServeHTTP(w, req)

	// Assertions básicas
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}
}
