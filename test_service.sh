#!/bin/bash

# Script para probar el servicio de autenticación
echo "🧪 Probando it-auth-service..."

BASE_URL="http://localhost:8082"

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para hacer peticiones HTTP
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    
    echo -e "${YELLOW}Testing ${method} ${endpoint}${NC}"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -d "$data" \
            "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} "${BASE_URL}${endpoint}")
    fi
    
    # Separar body y status code
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n 1)
    
    if [ "$status" = "$expected_status" ]; then
        echo -e "${GREEN}✅ SUCCESS: Status $status${NC}"
        echo "Response: $body"
    else
        echo -e "${RED}❌ FAILED: Expected $expected_status, got $status${NC}"
        echo "Response: $body"
    fi
    echo "---"
}

# Verificar que el servicio esté corriendo
echo "🔍 Verificando que el servicio esté corriendo..."
if ! curl -s "${BASE_URL}/api/v1/health" > /dev/null; then
    echo -e "${RED}❌ El servicio no está corriendo en ${BASE_URL}${NC}"
    echo "Asegúrate de ejecutar: ./scripts/run-local.sh"
    exit 1
fi

echo -e "${GREEN}✅ Servicio está corriendo${NC}"
echo ""

# Probar endpoints
echo "🚀 Probando endpoints..."
echo ""

# Health check
test_endpoint "GET" "/api/v1/health" "" "200"

# Readiness check
test_endpoint "GET" "/api/v1/ready" "" "200"

# API Info
test_endpoint "GET" "/api/v1/info" "" "200"

# Test Firebase login (sin token válido, debería fallar)
test_endpoint "POST" "/api/v1/auth/firebase-login" '{"firebase_token":"invalid_token","provider":"google"}' "400"

# Test Firebase register (sin token válido, debería fallar)
test_endpoint "POST" "/api/v1/auth/firebase-register" '{"firebase_token":"invalid_token","provider":"google","registration_data":{"username":"testuser"}}' "400"

# Test refresh token (sin token válido, debería fallar)
test_endpoint "POST" "/api/v1/auth/refresh-token" '{"firebase_token":"invalid_token"}' "400"

# Test logout (sin autorización, debería fallar)
test_endpoint "POST" "/api/v1/auth/logout" '{}' "401"

# Test get profile (sin autorización, debería fallar)
test_endpoint "GET" "/api/v1/users/profile" "" "401"

# Test update profile (sin autorización, debería fallar)
test_endpoint "PUT" "/api/v1/users/profile" '{"first_name":"Test"}' "401"

# Test list users (sin autorización, debería fallar)
test_endpoint "GET" "/api/v1/users" "" "401"

echo ""
echo -e "${GREEN}🎉 Pruebas completadas!${NC}"
echo ""
echo "📝 Notas:"
echo "- Los endpoints de autenticación fallan como se esperaba (tokens inválidos)"
echo "- Los endpoints protegidos fallan como se esperaba (sin autorización)"
echo "- Para probar con tokens reales, necesitarás tokens de Firebase válidos"
echo ""
echo "🔗 URLs útiles:"
echo "- Health: ${BASE_URL}/api/v1/health"
echo "- API Info: ${BASE_URL}/api/v1/info"
echo "- Swagger (si está habilitado): ${BASE_URL}/swagger/index.html"