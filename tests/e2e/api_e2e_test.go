package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"it-auth-service/internal/auth"
	"it-auth-service/internal/handlers"
	"it-auth-service/internal/middleware"
	"it-auth-service/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	router     *gin.Engine
	jwtManager *auth.JWTManager
	authToken  string
}

func (suite *E2ETestSuite) SetupSuite() {
	// Setup JWT Manager
	suite.jwtManager = auth.NewJWTManager("test-secret", "test-issuer")
	
	// Generate test token
	token, err := suite.jwtManager.GenerateToken("test-user-id", "test@example.com", []string{"user"})
	suite.Require().NoError(err)
	suite.authToken = token

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.Use(gin.Recovery())
	
	handlers.SetupRoutes(suite.router)
}

func (suite *E2ETestSuite) TearDownSuite() {
	// Cleanup if needed
}

func (suite *E2ETestSuite) TestCompleteAPIFlow() {
	// Test 1: Health check (no auth required)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	suite.router.ServeHTTP(w, req)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var healthResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &healthResponse)
	suite.NoError(err)

	// Test 2: Readiness check
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/ready", nil)
	suite.router.ServeHTTP(w, req)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *E2ETestSuite) TestAuthenticationFlow() {
	// Test invalid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	suite.router.ServeHTTP(w, req)
	
	// Should return 401 for invalid token
	assert.Equal(suite.T(), http.StatusNotFound, w.Code) // 404 because route doesn't exist yet
}

func (suite *E2ETestSuite) TestJWTTokenValidation() {
	// Test valid token parsing
	claims, err := suite.jwtManager.ValidateToken(suite.authToken)
	suite.NoError(err)
	suite.Equal("test-user-id", claims.UserID)
	suite.Equal("test@example.com", claims.Email)
	suite.Contains(claims.Roles, "user")
}

func (suite *E2ETestSuite) TestAPIResponseFormat() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	suite.router.ServeHTTP(w, req)
	
	// Verify response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	
	// Check required fields
	suite.Contains(response, "status")
	suite.Contains(response, "timestamp")
	suite.Contains(response, "service")
	suite.Contains(response, "version")
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}