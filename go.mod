module it-auth-service

go 1.21

require (
	firebase.google.com/go/v4 v4.12.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.27.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/hashicorp/vault/api v1.10.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.26.0
	golang.org/x/time v0.5.0
	google.golang.org/api v0.126.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
	github.com/prometheus/client_golang v1.17.0
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/swag v1.16.2
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/testcontainers/testcontainers-go v0.26.0
	github.com/golang-migrate/migrate/v4 v4.17.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	github.com/redis/go-redis/v9 v9.3.0
)