/**
 * k6 — Auth Service: Flujo completo de autenticación
 * =====================================================================
 * Flujo ejecutado en orden secuencial:
 *   1. POST /register   → Crea un usuario único (con UUID en email/username)
 *   2. POST /login      → Inicia sesión, obtiene access_token + refresh_token
 *   3. GET  /me         → Inspecciona los claims del JWT (usuario autenticado)
 *   4. POST /logout     → Revoca el JWT (lo manda a la blacklist) y el refresh_token
 *   5. GET  /me         → Intenta usar el JWT revocado → debe responder 401
 *   6. POST /refresh    → Intenta usar el refresh_token revocado → debe responder 401
 *
 * Uso:
 *   k6 run tests/k6/auth_flow.js
 *
 * Con reporte visual:
 *   k6 run --out json=tests/k6/results.json tests/k6/auth_flow.js
 * =====================================================================
 */

import http from "k6/http";
import { check, group, sleep } from "k6";
import { uuidv4 } from "./lib/k6-utils.js";

// ─── CONFIGURACIÓN ──────────────────────────────────────────────────────────

const BASE_URL = __ENV.BASE_URL || "http://localhost:9090";

export const options = {
  // Escenario: 1 VU, ejecuta el flujo 1 vez — test de contrato funcional
  vus: 1,
  iterations: 1,

  // Umbrales de éxito globales
  thresholds: {
    // Todas las peticiones deben tener tasa de error < 1%
    http_req_failed: ["rate<0.01"],
    // El flujo completo no debe superar 5s
    http_req_duration: ["p(95)<5000"],
    // Todos los checks deben pasar al 100%
    checks: ["rate==1.0"],
  },
};

// ─── HEADERS COMUNES ─────────────────────────────────────────────────────────

const JSON_HEADERS = { "Content-Type": "application/json" };

function authHeaders(token) {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

/**
 * logCurl — Imprime el equivalente curl de la peticion para reproducirla manualmente.
 * Util para depuracion y para generar reportes de monitoreo.
 */
function logCurl(method, url, headers, body) {
  let cmd = `curl -s -X ${method}`;

  Object.keys(headers).forEach(function (key) {
    // Ocultar el token JWT completo en el log — mostrar solo los primeros 30 chars
    let val = headers[key];
    if (key === "Authorization" && val.length > 40) {
      val = val.substring(0, 37) + "...";
    }
    cmd += ` \\\n    -H '${key}: ${val}'`;
  });

  if (body) {
    // Ocultar passwords en el log
    let safeBody = body;
    try {
      if (body !== "{}") {
        const parsed = JSON.parse(body);
        if (parsed.password) parsed.password = "[REDACTED]";
        safeBody = JSON.stringify(parsed);
      }
    } catch (_) {}
    cmd += ` \\\n    -d '${safeBody}'`;
  }

  cmd += ` \\\n    '${url}'`;
  console.log(`\n[CURL] ${cmd}`);
}

/**
 * logResponse — Imprime status + body parseado de la respuesta.
 */
function logResponse(label, res) {
  console.log(`\n[RESPONSE] ${label}`);
  console.log(`  Status   : ${res.status}`);
  console.log(`  Duration : ${res.timings.duration.toFixed(2)}ms`);
  try {
    const body = JSON.parse(res.body);
    console.log(`  Body     : ${JSON.stringify(body, null, 2)}`);
  } catch (_) {
    console.log(`  Body     : ${res.body}`);
  }
}

/**
 * logRequest — Wrapper que combina logCurl + logResponse en una sola llamada.
 * Centraliza todo el trazado de una peticion HTTP en el log.
 */
function logRequest(label, method, url, headers, body, res) {
  const divider = "=".repeat(60);
  console.log(`\n${divider}`);
  console.log(`[REQUEST] ${label}`);
  console.log(divider);
  logCurl(method, url, headers, body);
  logResponse(label, res);
}

/** Genera datos de usuario unicos para que el test sea idempotente */
function generateUser() {
  const id = uuidv4().substring(0, 8);
  return {
    username: `testuser_${id}`,
    email: `testuser_${id}@k6.homelab.dev`,
    password: "K6TestPass123!",
  };
}

// ─── ESCENARIO PRINCIPAL ─────────────────────────────────────────────────────

export default function () {
  const user = generateUser();

  console.log("\n╔══════════════════════════════════════════════════════╗");
  console.log("║    AUTH SERVICE — Flujo Completo de Autenticación   ║");
  console.log("╚══════════════════════════════════════════════════════╝");
  console.log(`  Usuario de prueba: ${user.username} | ${user.email}`);

  let accessToken = "";
  let refreshToken = "";

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 1: REGISTRO DE USUARIO
  // ═══════════════════════════════════════════════════════════════════════
  group("1. Registro de usuario", function () {
    const payload = JSON.stringify({
      username: user.username,
      email: user.email,
      password: user.password,
    });

    const res = http.post(`${BASE_URL}/api/v1/auth/register`, payload, {
      headers: JSON_HEADERS,
    });

    logRequest("POST /api/v1/auth/register", "POST",
      `${BASE_URL}/api/v1/auth/register`, JSON_HEADERS, payload, res);

    check(res, {
      "✅ [REGISTER] HTTP 201 Created": (r) => r.status === 201,
      "✅ [REGISTER] success=true": (r) => JSON.parse(r.body).success === true,
      "✅ [REGISTER] data.id existe": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.id !== undefined && b.data.id !== "";
      },
      "✅ [REGISTER] data.username coincide": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.username === user.username;
      },
      "✅ [REGISTER] data.email coincide": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.email === user.email;
      },
      "✅ [REGISTER] password_hash NO expuesto": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.password_hash === undefined;
      },
    });

    sleep(0.3);
  });

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 2: LOGIN — Obtener access_token + refresh_token
  // ═══════════════════════════════════════════════════════════════════════
  group("2. Login — obtener tokens", function () {
    const payload = JSON.stringify({
      email: user.email,
      password: user.password,
    });

    const res = http.post(`${BASE_URL}/api/v1/auth/login`, payload, {
      headers: JSON_HEADERS,
    });

    logRequest("POST /api/v1/auth/login", "POST",
      `${BASE_URL}/api/v1/auth/login`, JSON_HEADERS, payload, res);

    check(res, {
      "✅ [LOGIN] HTTP 200 OK": (r) => r.status === 200,
      "✅ [LOGIN] success=true": (r) => JSON.parse(r.body).success === true,
      "✅ [LOGIN] data.token presente": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.token && b.data.token.length > 20;
      },
      "✅ [LOGIN] cookie refresh_token presente y segura": (r) => {
        return (
          r.cookies &&
          r.cookies.refresh_token &&
          r.cookies.refresh_token.length > 0 &&
          r.cookies.refresh_token[0].value.length > 20
        );
      },
      "✅ [LOGIN] token tiene formato JWT (3 segmentos)": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.token.split(".").length === 3;
      },
    });

    // Extraer tokens para los siguientes pasos
    if (res.status === 200) {
      const body = JSON.parse(res.body);
      accessToken = body.data.token;
      
      if (res.cookies && res.cookies.refresh_token && res.cookies.refresh_token.length > 0) {
        refreshToken = res.cookies.refresh_token[0].value;
      }
      
      console.log(`  ▶ access_token  : ${accessToken.substring(0, 40)}...`);
      console.log(`  ▶ refresh_token : ${refreshToken.substring(0, 30)}... (via Cookie)`);
    }

    sleep(0.3);
  });

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 3: GET /me — Inspeccionar claims del JWT
  // ═══════════════════════════════════════════════════════════════════════
  group("3. GET /me — Inspeccionar claims del JWT", function () {
    if (!accessToken) {
      console.error("  ✗ Sin access_token, saltando /me");
      return;
    }

    const res = http.get(`${BASE_URL}/api/v1/auth/me`, {
      headers: authHeaders(accessToken),
    });

    logRequest("GET /api/v1/auth/me (token valido)", "GET",
      `${BASE_URL}/api/v1/auth/me`, authHeaders(accessToken), null, res);

    check(res, {
      "✅ [ME] HTTP 200 OK": (r) => r.status === 200,
      "✅ [ME] success=true": (r) => JSON.parse(r.body).success === true,
      "✅ [ME] data.user_id presente": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.user_id && b.data.user_id.length > 0;
      },
      "✅ [ME] data.username coincide": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.username === user.username;
      },
      "✅ [ME] data.email coincide": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.email === user.email;
      },
      "✅ [ME] data.roles es un array": (r) => {
        const b = JSON.parse(r.body);
        return b.data && Array.isArray(b.data.roles);
      },
    });

    // Mostrar los claims completos
    if (res.status === 200) {
      const claims = JSON.parse(res.body).data;
      console.log("\n  ── Claims del JWT ──────────────────────────────");
      console.log(`  user_id  : ${claims.user_id}`);
      console.log(`  username : ${claims.username}`);
      console.log(`  email    : ${claims.email}`);
      console.log(`  roles    : ${JSON.stringify(claims.roles)}`);
    }

    sleep(0.3);
  });

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 4: LOGOUT — Revocar JWT y refresh_token
  // ═══════════════════════════════════════════════════════════════════════
  group("4. Logout — revocar tokens", function () {
    if (!accessToken) {
      console.error("  ✗ Sin access_token, saltando logout");
      return;
    }

    const payload = JSON.stringify({});

    const headers = authHeaders(accessToken);
    // Inyectar cookie manualmente para simular el navegador
    headers["Cookie"] = `refresh_token=${refreshToken}`;

    const res = http.post(`${BASE_URL}/api/v1/auth/logout`, payload, {
      headers: headers,
    });

    logRequest("POST /api/v1/auth/logout", "POST",
      `${BASE_URL}/api/v1/auth/logout`, headers, payload, res);

    check(res, {
      "✅ [LOGOUT] HTTP 200 OK": (r) => r.status === 200,
      "✅ [LOGOUT] success=true": (r) => JSON.parse(r.body).success === true,
      "✅ [LOGOUT] mensaje de confirmación": (r) => {
        const b = JSON.parse(r.body);
        return b.data && b.data.message !== undefined;
      },
    });

    sleep(0.5); // Pequeña pausa para que Redis procese la blacklist
  });

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 5: GET /me con JWT revocado — debe devolver 401
  // ═══════════════════════════════════════════════════════════════════════
  group("5. GET /me con JWT en blacklist → 401 esperado", function () {
    if (!accessToken) {
      console.error("  ✗ Sin access_token, saltando validación blacklist");
      return;
    }

    const res = http.get(`${BASE_URL}/api/v1/auth/me`, {
      headers: authHeaders(accessToken),
    });

    logRequest("GET /api/v1/auth/me (JWT en blacklist)", "GET",
      `${BASE_URL}/api/v1/auth/me`, authHeaders(accessToken), null, res);

    check(res, {
      "✅ [BLACKLIST-JWT] HTTP 401 Unauthorized": (r) => r.status === 401,
      "✅ [BLACKLIST-JWT] success=false": (r) =>
        JSON.parse(r.body).success === false,
      "✅ [BLACKLIST-JWT] error.code presente": (r) => {
        const b = JSON.parse(r.body);
        return b.error && b.error.code !== undefined;
      },
    });

    sleep(0.3);
  });

  // ═══════════════════════════════════════════════════════════════════════
  // PASO 6: POST /refresh con refresh_token revocado — debe devolver 401
  // ═══════════════════════════════════════════════════════════════════════
  group(
    "6. POST /refresh con refresh_token revocado → 401 esperado",
    function () {
      if (!refreshToken) {
        console.error("  ✗ Sin refresh_token, saltando validación");
        return;
      }

      const payload = JSON.stringify({});

      const headers = Object.assign({}, JSON_HEADERS);
      // Inyectar cookie manualmente para simular el navegador
      headers["Cookie"] = `refresh_token=${refreshToken}`;

      const res = http.post(`${BASE_URL}/api/v1/auth/refresh`, payload, {
        headers: headers,
      });

      logRequest("POST /api/v1/auth/refresh (refresh_token revocado)", "POST",
        `${BASE_URL}/api/v1/auth/refresh`, headers, payload, res);

      check(res, {
        "✅ [BLACKLIST-RT] HTTP 401 Unauthorized": (r) => r.status === 401,
        "✅ [BLACKLIST-RT] success=false": (r) =>
          JSON.parse(r.body).success === false,
        "✅ [BLACKLIST-RT] error.code presente": (r) => {
          const b = JSON.parse(r.body);
          return b.error && b.error.code !== undefined;
        },
      });
    }
  );

  // ─── RESUMEN FINAL ────────────────────────────────────────────────────
  console.log("\n╔══════════════════════════════════════════════════════╗");
  console.log("║              Flujo completado ✓                     ║");
  console.log("╚══════════════════════════════════════════════════════╝");
}
