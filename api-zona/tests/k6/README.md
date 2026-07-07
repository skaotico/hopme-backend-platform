# k6 — Auth Service Tests

Tests de contrato y carga para el servicio de autenticación. Corren en Docker, **no requieren k6 instalado** en la máquina del desarrollador.

---

## Estructura

```
tests/k6/
├── auth_flow.js        # Script principal — flujo completo de autenticación
├── docker-compose.yml  # Compose para levantar k6 como servicio
├── run-k6.sh           # Script shell de conveniencia (Linux / Mac / WSL)
├── lib/
│   └── k6-utils.js     # Dependencia local vendorizada (uuidv4)
└── results/            # Reportes JSON generados (ignorados por git)
```

---

## Requisitos

- Docker Desktop corriendo
- El servicio de auth levantado (ver instrucciones en el README raíz)

---

## Uso rápido

### Opción A — Shell script (recomendado)

```sh
# Dar permisos de ejecución (solo la primera vez)
chmod +x tests/k6/run-k6.sh

# Flujo completo contra localhost:9090 (por defecto)
./tests/k6/run-k6.sh

# Apuntando a otra URL
./tests/k6/run-k6.sh --url http://192.168.1.20:9090

# Con reporte JSON guardado en tests/k6/results/
./tests/k6/run-k6.sh --report

# Combinado
./tests/k6/run-k6.sh --url http://192.168.1.20:9090 --report
```

> En **Windows** ejecutar desde WSL, Git Bash o cualquier terminal con `sh`.

### Opción B — Docker Compose

```bash
# Desde la raíz del proyecto
docker compose -f tests/k6/docker-compose.yml run --rm k6

# Con URL personalizada
BASE_URL=http://192.168.1.20:9090 docker compose -f tests/k6/docker-compose.yml run --rm k6
```

### Opción C — Docker directo

```bash
docker run --rm \
  -v ./tests/k6:/scripts:ro \
  -e BASE_URL=http://host.docker.internal:9090 \
  --add-host host.docker.internal:host-gateway \
  grafana/k6:latest run /scripts/auth_flow.js
```

---

## Flujo que ejecuta `auth_flow.js`

| # | Endpoint | Qué valida |
|---|----------|------------|
| 1 | `POST /register` | Crea usuario único, HTTP 201, `password_hash` no expuesto |
| 2 | `POST /login` | HTTP 200, JWT con 3 segmentos, refresh_token presente |
| 3 | `GET /me` | HTTP 200, claims: `user_id`, `username`, `email`, `roles[]` |
| 4 | `POST /logout` | HTTP 200, revoca JWT (blacklist Redis) + refresh_token (DB) |
| 5 | `GET /me` | **HTTP 401** — JWT ya en blacklist, rechazado |
| 6 | `POST /refresh` | **HTTP 401** — refresh_token revocado en PostgreSQL |

> El usuario creado en cada ejecución es único (UUID en email/username), por lo que el test es **idempotente** y se puede repetir sin limpiar la base de datos.

---

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `BASE_URL` | `http://host.docker.internal:9090` | URL base del servicio de auth |

> En Windows/Mac con Docker Desktop, `host.docker.internal` resuelve automáticamente al `localhost` del host. En Linux, el compose ya incluye `extra_hosts: host-gateway`.

---

## Interpretar resultados

Una ejecución exitosa muestra todos los checks en verde:

```
✓ checks.........................: 100.00% ✓ 18  ✗ 0
  http_req_duration..............: avg=45ms  p(95)=120ms
  http_req_failed................: 0.00%   ✓ 0   ✗ 6
```

Si algún check falla, k6 termina con **exit code 1** y el paso concreto aparece marcado con `✗`.
