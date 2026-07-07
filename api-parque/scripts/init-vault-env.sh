#!/bin/bash

# Leer token de parámetro o variable de entorno
VAULT_TOKEN=${1:-$VAULT_TOKEN}

if [ -z "$VAULT_TOKEN" ]; then
    echo -n "Ingresa tu Vault Token (o define la variable de entorno VAULT_TOKEN): "
    read -r VAULT_TOKEN
fi

if [ -z "$VAULT_TOKEN" ]; then
    echo -e "\033[0;31mERROR: El Vault Token es requerido para continuar.\033[0m"
    exit 1
fi

VAULT_ADDR="http://192.168.1.20:8200"
VAULT_SECRET_PATH="v1/retromarket/data/dev/api-parque"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

echo "Obteniendo secretos desde Vault: $VAULT_ADDR/$VAULT_SECRET_PATH"

RESPONSE=$(curl -s -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/$VAULT_SECRET_PATH")

# Verificar si obtuvimos una respuesta exitosa
if echo "$RESPONSE" | grep -q '"data"'; then
    echo "# ==============================================================================" > "$ENV_FILE"
    echo "# ARCHIVO .ENV GENERADO AUTOMÁTICAMENTE DESDE HASHICORP VAULT" >> "$ENV_FILE"
    echo "# Generado el: $(date)" >> "$ENV_FILE"
    echo "# ==============================================================================" >> "$ENV_FILE"
    echo "" >> "$ENV_FILE"

    # Intentar parsear el JSON usando jq (preferido) o python3 (fallback en macOS)
    if command -v jq &> /dev/null; then
        echo "$RESPONSE" | jq -r '.data.data | to_entries | .[] | "\(.key)=\(.value)"' >> "$ENV_FILE"
    elif command -v python3 &> /dev/null; then
        echo "Advertencia: 'jq' no está instalado. Utilizando Python 3 como alternativa..."
        echo "$RESPONSE" | python3 -c "
import sys, json
try:
    res = json.load(sys.stdin)
    secrets = res['data']['data']
    for k, v in secrets.items():
        print(f'{k}={v}')
except Exception as e:
    print(f'Error al procesar JSON: {e}', file=sys.stderr)
    sys.exit(1)
" >> "$ENV_FILE"
    else
        echo -e "\033[0;31mERROR: Se requiere 'jq' o 'python3' en el sistema para parsear los secretos JSON.\033[0m"
        exit 1
    fi

    echo -e "\033[0;32mArchivo .env generado exitosamente en: $ENV_FILE\033[0m"
else
    echo -e "\033[0;31mERROR: No se pudieron obtener los secretos de Vault.\033[0m"
    echo "Respuesta de Vault: $RESPONSE"
    exit 1
fi
