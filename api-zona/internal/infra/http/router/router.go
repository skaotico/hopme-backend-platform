package router

import (
	"net/http"

	zonaHTTP "c4-zona/internal/infra/http"
	"c4-zona/internal/domain/port"
)

// NewRouter crea y configura un multiplexor de peticiones
func NewRouter(
	uc port.ZonaUseCase,
	loggingMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Healthcheck
	mux.Handle("GET /api/v1/health", zonaHTTP.HealthHandler())

	// CRUD Zonas
	mux.Handle("POST /api/v1/zonas", zonaHTTP.CreateZonaHandler(uc))
	mux.Handle("GET /api/v1/zonas", zonaHTTP.ListZonasHandler(uc))
	mux.Handle("GET /api/v1/zonas/{id}", zonaHTTP.GetZonaHandler(uc))
	mux.Handle("PUT /api/v1/zonas/{id}", zonaHTTP.UpdateZonaHandler(uc))
	mux.Handle("DELETE /api/v1/zonas/{id}", zonaHTTP.DeleteZonaHandler(uc))

	mux.Handle("GET /swagger/", zonaHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", zonaHTTP.SwaggerJSONHandler())

	var handler http.Handler = mux
	if loggingMiddleware != nil {
		handler = loggingMiddleware(handler)
	}

	return handler
}
