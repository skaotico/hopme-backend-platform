#!/usr/bin/env sh
# =============================================================================
# run-k6.sh - Ejecutor de tests k6 en Docker con logging completo
#
# Genera en tests/k6/results/ por cada ejecucion:
#   run_TIMESTAMP.log         -> Log completo (consola + requests + checks)
#   metrics_TIMESTAMP.json    -> Metricas en formato JSON (para monitoreo)
#
# Uso:
#   ./tests/k6/run-k6.sh
#   ./tests/k6/run-k6.sh --url http://192.168.1.20:9090
#   ./tests/k6/run-k6.sh --script /scripts/otro.js
#
# Requiere: Docker corriendo
# =============================================================================

set -e

# -- Defaults ------------------------------------------------------------------
BASE_URL="http://host.docker.internal:9090"
SCRIPT="/scripts/auth_flow.js"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RESULTS_DIR="${SCRIPT_DIR}/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

LOG_FILE="${RESULTS_DIR}/run_${TIMESTAMP}.log"
METRICS_FILE="/scripts/results/metrics_${TIMESTAMP}.json"
METRICS_FILE_HOST="${RESULTS_DIR}/metrics_${TIMESTAMP}.json"

# -- Parsear argumentos --------------------------------------------------------
while [ $# -gt 0 ]; do
  case "$1" in
    --url)    BASE_URL="$2"; shift 2 ;;
    --script) SCRIPT="$2";   shift 2 ;;
    *)
      echo "Opcion desconocida: $1"
      echo "Uso: $0 [--url URL] [--script PATH]"
      exit 1
      ;;
  esac
done

# -- Crear carpeta de resultados -----------------------------------------------
mkdir -p "$RESULTS_DIR"

# -- Banner en consola y log ---------------------------------------------------
BANNER=$(cat <<EOF

======================================================
  k6 - Auth Service Test Runner (Docker)
======================================================
  Timestamp : ${TIMESTAMP}
  BASE_URL  : ${BASE_URL}
  Script    : ${SCRIPT}
  Log       : results/run_${TIMESTAMP}.log
  Metrics   : results/metrics_${TIMESTAMP}.json
======================================================
EOF
)

echo "$BANNER"
echo "$BANNER" > "$LOG_FILE"

# -- Ejecutar k6 en Docker -----------------------------------------------------
# Flags importantes:
#   --http-debug=full  -> imprime headers y body de cada request/response
#   --out json=FILE    -> exporta todas las metricas en JSON para monitoreo
#   K6_LOG_FORMAT=json -> logs internos de k6 en JSON estructurado
#
# El volumen NO es :ro para poder escribir los resultados desde dentro del contenedor.
# tee duplica la salida: pantalla + archivo de log simultaneamente.
#
# shellcheck disable=SC2086
docker run --rm \
  -v "${SCRIPT_DIR}:/scripts" \
  -e "BASE_URL=${BASE_URL}" \
  --add-host "host.docker.internal:host-gateway" \
  grafana/k6:latest \
  run \
    --out "json=${METRICS_FILE}" \
    ${SCRIPT} \
  2>&1 | tee -a "$LOG_FILE"

EXIT_CODE=${PIPESTATUS:-$?}

# -- Resumen final -------------------------------------------------------------
SUMMARY=$(cat <<EOF

======================================================
  Resultado de la ejecucion
======================================================
  Status  : $([ "$EXIT_CODE" = "0" ] && echo "OK - Todos los checks pasaron" || echo "FAIL - Algunos checks fallaron")
  Log     : ${LOG_FILE}
  Metrics : ${METRICS_FILE_HOST}
======================================================
EOF
)

echo "$SUMMARY"
echo "$SUMMARY" >> "$LOG_FILE"

exit "$EXIT_CODE"
