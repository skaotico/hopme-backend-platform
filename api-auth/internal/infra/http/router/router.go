package router

import (
	"net/http"
	"time"

	authHTTP "c4-auth/internal/infra/http"
	"c4-auth/internal/infra/http/middleware"
	"c4-auth/internal/domain/port"
)

// NewRouter crea y configura un multiplexor de peticiones nativo para toda la API de autenticación
func NewRouter(
	registerUC port.RegisterUseCase,
	loginUC port.LoginUseCase,
	refreshUC port.RefreshTokenUseCase,
	logoutUC port.LogoutUseCase,
	jwtMiddleware func(http.Handler) http.Handler,
	loggingMiddleware func(http.Handler) http.Handler,
	timeoutRegister time.Duration,
	timeoutLogin time.Duration,
	timeoutMe time.Duration,
	cookieSecure bool,
	rtExpiration time.Duration,
) http.Handler {
	mux := http.NewServeMux()

	// Endpoint público de Healthcheck
	mux.Handle("GET /api/v1/health", authHTTP.HealthHandler())

	// Aplicar timeout al endpoint de registro
	mux.Handle("POST /api/v1/auth/register", middleware.Timeout(timeoutRegister, "Registro")(authHTTP.RegisterHandler(registerUC)))
	
	// Aplicar timeout al endpoint de login
	mux.Handle("POST /api/v1/auth/login", middleware.Timeout(timeoutLogin, "Login")(authHTTP.LoginHandler(loginUC, cookieSecure, rtExpiration)))

	// Endpoint para refrescar tokens
	mux.Handle("POST /api/v1/auth/refresh", middleware.Timeout(timeoutLogin, "Refresh")(authHTTP.RefreshHandler(refreshUC, cookieSecure, rtExpiration)))

	mux.Handle("GET /swagger/", authHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", authHTTP.SwaggerJSONHandler())

	// Aplicar timeout y verificación JWT al endpoint protegido
	meHandler := authHTTP.MeHandler()
	mux.Handle("GET /api/v1/auth/me", jwtMiddleware(middleware.Timeout(timeoutMe, "Me")(meHandler)))

	// Endpoint para cerrar sesión (requiere JWT para añadirlo a blacklist y opcionalmente el refresh_token en cookie)
	mux.Handle("POST /api/v1/auth/logout", jwtMiddleware(middleware.Timeout(timeoutMe, "Logout")(authHTTP.LogoutHandler(logoutUC, cookieSecure))))

	var handler http.Handler = mux
	handler = loggingMiddleware(handler)

	return handler
}
