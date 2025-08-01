# 🔐 IT Auth Service - Guía de Testing Local

Esta guía te ayudará a configurar y probar el servicio de autenticación IT Auth Service localmente con Firebase.

## 📋 Requisitos Previos

- Go 1.22+
- PostgreSQL (configurado según tu `.env`)
- Firebase Project configurado
- Postman (para testing de API)

## 🚀 Configuración Inicial

### 1. **Variables de Entorno**

Tu archivo `.env` ya está configurado:

```env
# Auth Service Configuration
PORT=8082
ENVIRONMENT=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=35.227.10.150
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=p?<MJap]Lqm]LO6G
DB_NAME=it_db_chatbot

# Firebase Configuration
FIREBASE_PROJECT_ID=innovatech-agc

# Security
JWT_SECRET=9c8fc6c3cc8ecc3190e2a7bf5e2e8463b487be2189b2a91ad543fd3018cbcf8c
```

### 2. **Configuración de Firebase**

Necesitas crear el archivo `firebase-service-account.json` con las credenciales de tu proyecto Firebase:

```bash
# Copia el archivo de ejemplo
cp firebase-service-account.example.json firebase-service-account.json

# Edita el archivo con tus credenciales reales de Firebase
nano firebase-service-account.json
```

**Estructura del archivo:**
```json
{
  "type": "service_account",
  "project_id": "innovatech-agc",
  "private_key_id": "tu-private-key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\nTU_PRIVATE_KEY_AQUI\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xxxxx@innovatech-agc.iam.gserviceaccount.com",
  "client_id": "tu-client-id",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-xxxxx%40innovatech-agc.iam.gserviceaccount.com"
}
```

### 3. **Instalación de Dependencias**

```bash
# Instalar dependencias de Go
go mod tidy
go mod download
```

## 🏃‍♂️ Ejecutar el Servicio

```bash
# Desde la raíz del proyecto
go run cmd/main.go
```

El servicio estará disponible en: `http://localhost:8082`

## 📊 Endpoints Disponibles

### 🏥 **Health Checks**
- `GET /api/v1/health` - Estado del servicio
- `GET /api/v1/ready` - Verificación de disponibilidad
- `GET /api/v1/info` - Información de la API

### 🔥 **Firebase Authentication**
- `POST /api/v1/auth/firebase-login` - Login con token Firebase
- `POST /api/v1/auth/firebase-register` - Registro con token Firebase
- `POST /api/v1/auth/refresh-token` - Renovar token JWT

### 👥 **User Management**
- `GET /api/v1/users/profile` - Obtener perfil de usuario
- `PUT /api/v1/users/profile` - Actualizar perfil
- `GET /api/v1/users` - Listar usuarios (Admin)

## 🧪 Testing con Postman

### 1. **Importar Colección**

1. Abre Postman
2. Importa el archivo `postman_collection.json`
3. La colección incluye:
   - ✅ Tests automáticos con emojis
   - 🔄 Manejo automático de tokens
   - 📝 Variables de entorno preconfiguradas

### 2. **Variables de Entorno**

La colección incluye estas variables:

```json
{
  "base_url": "http://localhost:8082",
  "jwt_token": "",
  "user_id": "",
  "firebase_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.example_firebase_token_here"
}
```

### 3. **Flujo de Testing Recomendado**

#### **Paso 1: Verificar Health Checks**
1. 💚 **Health Check** - `GET /api/v1/health`
2. 🟢 **Readiness Check** - `GET /api/v1/ready`
3. 🔍 **API Info** - `GET /api/v1/info`

#### **Paso 2: Obtener Token de Firebase Real**

Para obtener un token de Firebase real, puedes:

**Opción A: Desde tu aplicación frontend**
```javascript
// En tu app frontend con Firebase
firebase.auth().currentUser.getIdToken(true)
  .then(token => {
    console.log('Firebase Token:', token);
    // Copia este token a Postman
  });
```

**Opción B: Usando Firebase Admin SDK (para testing)**
```bash
# Instalar Firebase CLI
npm install -g firebase-tools

# Login a Firebase
firebase login

# Generar token personalizado (requiere configuración adicional)
```

#### **Paso 3: Testing de Autenticación**

1. **🔥 Firebase Login**
   - Actualiza la variable `firebase_token` con un token real
   - Ejecuta el request
   - El token JWT se guardará automáticamente

2. **📝 Firebase Register**
   - Usa un token de Firebase de un usuario nuevo
   - Incluye datos adicionales en `registration_data`

3. **🔄 Refresh Token**
   - Usa un token de Firebase actualizado
   - Renueva tu JWT interno

#### **Paso 4: Testing de Gestión de Usuarios**

1. **👤 Get User Profile**
   - Requiere header `X-User-ID` (temporal para testing)
   - Usa el `user_id` guardado automáticamente

2. **✏️ Update User Profile**
   - Actualiza información del usuario
   - Campos permitidos: `username`, `first_name`, `last_name`, `photo_url`

3. **📋 List Users (Admin)**
   - Requiere header `X-Is-Admin: true` (temporal para testing)
   - Incluye paginación

## 🔧 Testing Sin Firebase (Desarrollo)

Para testing rápido sin Firebase, puedes:

### 1. **Probar Health Checks**
```bash
# Health check
curl http://localhost:8082/api/v1/health

# Readiness check
curl http://localhost:8082/api/v1/ready

# API info
curl http://localhost:8082/api/v1/info
```

### 2. **Simular Autenticación (Solo para desarrollo)**

Puedes modificar temporalmente el código para saltarse la verificación de Firebase:

```go
// En internal/services/firebase_auth.go
// Comentar temporalmente la verificación real y usar datos mock
```

## 🐛 Troubleshooting

### **Error: "Firebase credentials not found"**
- Verifica que `firebase-service-account.json` existe
- Verifica que `FIREBASE_PROJECT_ID` está configurado
- Verifica que las credenciales son válidas

### **Error: "Database connection failed"**
- Verifica la conexión a PostgreSQL
- Verifica las credenciales en `.env`
- Verifica que la base de datos `it_db_chatbot` existe

### **Error: "Invalid Firebase token"**
- Verifica que el token no ha expirado
- Verifica que el token es de tu proyecto Firebase
- Verifica que el usuario existe en Firebase

### **Error: "User not found"**
- Para endpoints de usuario, verifica el header `X-User-ID`
- Verifica que el usuario existe en la base de datos

## 📝 Logs y Debugging

El servicio incluye logging detallado:

```bash
# Ver logs en tiempo real
go run cmd/main.go | grep -E "(INFO|ERROR|WARN)"

# Ver solo errores
go run cmd/main.go 2>&1 | grep ERROR
```

## 🚀 Próximos Pasos

1. **Implementar middleware JWT** para autenticación real
2. **Agregar validación de roles** para endpoints de admin
3. **Implementar rate limiting**
4. **Agregar más tests unitarios**
5. **Configurar CI/CD** para despliegue automático

## 📞 Soporte

Si encuentras problemas:

1. Verifica los logs del servicio
2. Verifica la configuración de Firebase
3. Verifica la conexión a la base de datos
4. Revisa la documentación de Firebase Auth

¡Happy Testing! 🎉