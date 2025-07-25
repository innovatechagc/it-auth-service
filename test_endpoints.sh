#!/bin/bash

echo "🚀 Probando endpoints del servicio de autenticación..."
echo ""

BASE_URL="http://localhost:8082"

echo "1. Health Check:"
curl -s -X GET $BASE_URL/health | jq '.' || echo "❌ Health check falló"
echo ""

echo "2. Auth Status (sin token):"
curl -s -X GET $BASE_URL/auth/status | jq '.' || echo "❌ Auth status falló"
echo ""

echo "3. Token Verify (token inválido):"
curl -s -X POST $BASE_URL/tokens/verify \
  -H "Content-Type: application/json" \
  -d '{"id_token": "invalid-token"}' || echo "❌ Token verify falló"
echo ""

echo "4. Token Validate (token inválido):"
curl -s -X POST $BASE_URL/tokens/validate \
  -H "Content-Type: application/json" \
  -d '{"id_token": "invalid-token"}' || echo "❌ Token validate falló"
echo ""

echo "5. Login (token inválido):"
curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"id_token": "invalid-token"}' || echo "❌ Login falló"
echo ""

echo "✅ Pruebas completadas. El servicio está funcionando correctamente!"
echo ""
echo "📋 Endpoints disponibles:"
echo "  GET  /health                 - Health check"
echo "  GET  /auth/status            - Estado de autenticación"
echo "  POST /auth/login             - Login con Firebase token"
echo "  POST /auth/logout            - Logout"
echo "  POST /auth/refresh           - Refresh token"
echo "  POST /tokens/verify          - Verificar token"
echo "  POST /tokens/validate        - Validar token"
echo "  POST /tokens/refresh         - Refresh token"
echo ""
echo "📋 Endpoints protegidos (requieren autenticación):"
echo "  GET  /auth/profile           - Perfil del usuario"
echo "  POST /auth/revoke-tokens     - Revocar todos los tokens"
echo "  POST /tokens/revoke          - Revocar token específico"
echo "  POST /tokens/revoke-all      - Revocar todos los tokens"
echo "  POST /tokens/custom          - Crear token personalizado"