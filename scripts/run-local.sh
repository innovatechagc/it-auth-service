#!/bin/bash

# Script para ejecutar el servicio de autenticación localmente
echo "🚀 Iniciando it-auth-service en modo local..."

# Verificar que Go esté instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado. Por favor instala Go primero."
    exit 1
fi

# Cambiar al directorio del servicio de autenticación
cd "$(dirname "$0")/.."

# Cargar variables de entorno desde .env.local
if [ -f ".env.local" ]; then
    echo "📄 Cargando variables de entorno desde .env.local"
    export $(cat .env.local | grep -v '^#' | xargs)
else
    echo "⚠️  Archivo .env.local no encontrado. Usando variables de entorno del sistema."
fi

# Verificar que el archivo de credenciales de Firebase existe
if [ ! -f "$FIREBASE_SERVICE_ACCOUNT_PATH" ]; then
    echo "❌ Archivo de credenciales de Firebase no encontrado: $FIREBASE_SERVICE_ACCOUNT_PATH"
    echo "   Verifica que el archivo firebase-service-account.json esté presente."
    exit 1
fi

# Mostrar configuración
echo "🔧 Configuración:"
echo "   - DB_HOST: $DB_HOST"
echo "   - DB_PORT: $DB_PORT"
echo "   - DB_NAME: $DB_NAME"
echo "   - DB_USER: $DB_USER"
echo "   - PORT: $PORT"
echo "   - ENVIRONMENT: $ENVIRONMENT"
echo "   - LOG_LEVEL: $LOG_LEVEL"
echo "   - FIREBASE_PROJECT_ID: $FIREBASE_PROJECT_ID"
echo "   - FIREBASE_SERVICE_ACCOUNT_PATH: $FIREBASE_SERVICE_ACCOUNT_PATH"

# Verificar conexión a la base de datos
echo "🔍 Verificando conexión a la base de datos..."
if ! nc -z $DB_HOST $DB_PORT 2>/dev/null; then
    echo "❌ No se puede conectar a la base de datos en $DB_HOST:$DB_PORT"
    echo "   Verifica que la base de datos esté accesible desde tu red local."
    exit 1
fi

echo "✅ Conexión a la base de datos verificada"

# Descargar dependencias
echo "📦 Descargando dependencias..."
go mod download

# Compilar y ejecutar
echo "🔨 Compilando aplicación..."
go build -o bin/it-auth-service ./cmd

if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa"
    echo "🌟 Iniciando servidor en puerto $PORT..."
    echo "📍 Health check: http://localhost:$PORT/api/v1/health"
    echo "📍 API Base URL: http://localhost:$PORT/api/v1"
    echo ""
    echo "Presiona Ctrl+C para detener el servidor"
    echo "----------------------------------------"
    ./bin/it-auth-service
else
    echo "❌ Error en la compilación"
    exit 1
fi