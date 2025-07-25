package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func SetupRoutes(router *gin.Engine) {
	h := NewHandler()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", h.HealthCheck)
		api.GET("/ready", h.ReadinessCheck)
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Verifica el estado del servicio
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "it-auth-service",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}
	
	c.JSON(http.StatusOK, response)
}

// ReadinessCheck godoc
// @Summary Readiness check endpoint
// @Description Verifica si el servicio está listo para recibir tráfico
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *Handler) ReadinessCheck(c *gin.Context) {
	response := map[string]interface{}{
		"ready":     true,
		"service":   "it-auth-service",
		"timestamp": time.Now().UTC(),
	}
	
	c.JSON(http.StatusOK, response)
}

// Ejemplo de handler comentado para testing
/*
// GetExample godoc
// @Summary Get example data
// @Description Obtiene datos de ejemplo
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /example [get]
func (h *Handler) GetExample(c *gin.Context) {
	// Implementación de ejemplo
	c.JSON(http.StatusOK, gin.H{
		"message": "Example data",
		"data":    []string{"item1", "item2", "item3"},
	})
}

// CreateExample godoc
// @Summary Create example data
// @Description Crea datos de ejemplo
// @Tags example
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Example data"
// @Success 201 {object} map[string]interface{}
// @Router /example [post]
func (h *Handler) CreateExample(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Implementación de ejemplo
	c.JSON(http.StatusCreated, gin.H{
		"message": "Example created",
		"data":    request,
	})
}
*/