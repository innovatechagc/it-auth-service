# Servicio de Autenticación IT-Auth-Service

Microservicio de autenticación en Go que integra Firebase Auth, JWT y PostgreSQL. Diseñado para manejar autenticación, autorización y gestión de tokens de forma segura y escalable.

## 🚀 Características

- **Autenticación**: Firebase Auth integrado
- **JWT**: Manejo de tokens JWT personalizados
- **Base de datos**: PostgreSQL con GORM
- **Framework**: Gorilla Mux para HTTP routing
- **Logging**: Logrus para logging estructurado
- **Middleware**: Autenticación y CORS
- **Validación**: Validación de requests con go-playground/validator
- **Rate Limiting**: Control de tasa de requests
- **Testing**: Tests unitarios y de integración

## 📁 Estructura del Proyecto

```
├── cmd/                    # Comandos de la aplicación
├── internal/              # Código interno de la aplicación
│   ├── config/           # Configuración
│   ├── handlers/         # Handlers HTTP
│   ├── middleware/       # Middleware personalizado
│   └── services/         # Lógica de negocio
├── pkg/                  # Paquetes reutilizables
│   ├── logger/          # Logger personalizado
│   └── vault/           # Cliente de Vault
├── scripts/             # Scripts de inicialización
├── monitoring/          # Configuración de monitoreo
├── .env.*              # Archivos de configuración por entorno
├── docker-compose.yml  # Desarrollo local
├── Dockerfile         # Imagen de producción
└── Makefile          # Comandos de automatización
```

## 🛠️ Configuración Inicial

### 1. Clonar y configurar el proyecto

```bash
# Clonar el repositorio
git clone <repository-url>
cd it-auth-service

# Copiar configuración de ejemplo
cp .env.example .env

# Instalar dependencias
go mod tidy
```

### 2. Configurar variables de entorno

Edita `.env` con tus configuraciones:

```bash
# Auth Service Configuration
PORT=8082
ENVIRONMENT=development
LOG_LEVEL=debug

# Database Configuration (SHARED with it-app_user)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=itapp

# Firebase Configuration
FIREBASE_PROJECT_ID=innovatech-agc

# Rate Limiting
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200

# Security
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-2024
```

### 3. Configurar Firebase

**⚠️ IMPORTANTE: Nunca subas credenciales reales al repositorio**

1. Copia el archivo de ejemplo:
```bash
cp firebase-service-account.json.example firebase-service-account.json
```

2. Edita `firebase-service-account.json` con tus credenciales reales de Firebase Admin SDK
3. El archivo `firebase-service-account.json` está en `.gitignore` y NO se subirá al repositorio

**Para obtener las credenciales de Firebase:**
1. Ve a [Firebase Console](https://console.firebase.google.com/)
2. Selecciona tu proyecto
3. Ve a Configuración del proyecto > Cuentas de servicio
4. Genera una nueva clave privada
5. Descarga el archivo JSON y renómbralo a `firebase-service-account.json`

## 🚀 Desarrollo Local

### Opción 1: Ejecutar directamente

```bash
# Compilar y ejecutar
go build -o bin/auth-service ./cmd
./bin/auth-service

# O usando make
make build
make run
```

### Opción 2: Con Docker Compose (Recomendado)

```bash
# Levantar todos los servicios (app, postgres, vault, redis, prometheus)
make docker-dev

# Detener servicios
make docker-down
```

Servicios disponibles:
- **Auth API**: http://localhost:8082
- **PostgreSQL**: localhost:5432
- **Prometheus**: http://localhost:9090
- **Vault**: http://localhost:8200

### Probar el servicio

```bash
# Ejecutar script de pruebas
./test_endpoints.sh

# O probar manualmente
curl http://localhost:8082/health
```

## 🧪 Testing

```bash
# Ejecutar tests
make test

# Tests con cobertura
make test-coverage

# Tests con Docker
make docker-test

# Linting
make lint
```

## 📊 Endpoints Disponibles

### Health Checks
- `GET /health` - Estado del servicio

### Endpoints Públicos de Autenticación
- `POST /auth/login` - Login con Firebase ID token
- `POST /auth/logout` - Logout del usuario
- `GET /auth/status` - Verificar estado de autenticación
- `POST /auth/refresh` - Refresh token

### Endpoints de Tokens
- `POST /tokens/verify` - Verificar token de Firebase
- `POST /tokens/validate` - Validar token
- `POST /tokens/refresh` - Refresh token

### Endpoints Protegidos (requieren autenticación)
- `GET /auth/profile` - Obtener perfil del usuario
- `POST /auth/revoke-tokens` - Revocar todos los tokens del usuario
- `POST /tokens/revoke` - Revocar token específico
- `POST /tokens/revoke-all` - Revocar todos los tokens
- `POST /tokens/custom` - Crear token personalizado

### Ejemplos de Uso

```bash
# Health check
curl http://localhost:8082/health

# Verificar estado de autenticación (sin token)
curl http://localhost:8082/auth/status

# Login con Firebase token
curl -X POST http://localhost:8082/auth/login \
  -H "Content-Type: application/json" \
  -d '{"id_token": "your-firebase-id-token"}'

# Verificar token
curl -X POST http://localhost:8082/tokens/verify \
  -H "Content-Type: application/json" \
  -d '{"id_token": "your-firebase-id-token"}'

# Acceder a endpoint protegido
curl -X GET http://localhost:8082/auth/profile \
  -H "Authorization: Bearer your-firebase-id-token"
```

## 🔧 Configuración por Entornos

### Desarrollo Local
- Archivo: `.env.local`
- Base de datos: PostgreSQL local
- Vault: Opcional (comentado por defecto)
- Logs: Debug level

### Testing/QA
- Archivo: `.env.test`
- Base de datos: PostgreSQL de testing
- Vault: Instancia de testing
- Logs: Info level

### Producción
- Archivo: `.env.production`
- Variables desde GCP Secret Manager o Vault
- SSL requerido para BD
- Logs: Warn level

## 🐳 Docker

### Desarrollo
```bash
# Construir imagen
make docker-build

# Ejecutar contenedor
make docker-run
```

### Testing
```bash
# Ejecutar tests en contenedor
make docker-test
```

## ☁️ Despliegue en GCP Cloud Run

### Preparación
1. Configurar gcloud CLI
2. Habilitar Cloud Run API
3. Configurar Container Registry

### Deploy a Staging
```bash
# Build y push de imagen
docker build -t gcr.io/PROJECT_ID/microservice-template:latest .
docker push gcr.io/PROJECT_ID/microservice-template:latest

# Deploy
make deploy-staging
```

### Deploy a Producción
```bash
make deploy-prod
```

## 🔐 Manejo de Secretos

### Con Vault (Recomendado)
```go
// Ejemplo de uso
vaultClient, err := vault.NewClient(cfg.VaultConfig)
secrets, err := vaultClient.GetSecret("secret/myapp/database")
password := secrets["password"].(string)
```

### Variables de Entorno
Para desarrollo local, usar archivos `.env.*`

## 📈 Monitoreo y Métricas

### Métricas Disponibles
- `http_requests_total` - Total de requests HTTP
- `http_request_duration_seconds` - Duración de requests

### Prometheus
Configuración en `monitoring/prometheus.yml`

## 🔄 Personalización del Template

### 1. Cambiar nombre del módulo
Actualizar en `go.mod`:
```go
module github.com/company/tu-microservicio
```

### 2. Agregar nuevos endpoints
```go
// En internal/handlers/handlers.go
api.GET("/tu-endpoint", h.TuHandler)
```

### 3. Agregar servicios externos
```go
// En internal/services/
type ExternalService interface {
    CallAPI() error
}
```

### 4. Configurar base de datos
Descomentar y configurar en:
- `internal/config/config.go`
- Scripts de migración en `scripts/`

## 📝 Comandos Útiles

```bash
# Ver todos los comandos disponibles
make help

# Desarrollo
make deps          # Instalar dependencias
make build         # Compilar
make run           # Ejecutar
make test          # Tests
make lint          # Linting
make format        # Formatear código

# Docker
make docker-build  # Construir imagen
make docker-dev    # Entorno completo
make docker-test   # Tests en Docker

# Documentación
make swagger       # Generar docs Swagger
```

## 🤝 Contribución

1. Fork el proyecto
2. Crear feature branch (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push al branch (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 🆘 Soporte

Para preguntas o problemas:
1. Revisar la documentación
2. Buscar en issues existentes
3. Crear nuevo issue con detalles del problema

---

**Nota**: Este template incluye ejemplos comentados para facilitar el desarrollo. Descomenta y configura según las necesidades de tu microservicio.