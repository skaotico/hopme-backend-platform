package router

import (
	"net/http"

	parqueHTTP "c4-parque/internal/infra/http"
	"c4-parque/internal/domain/port"
)

// NewRouter crea y configura un multiplexor de peticiones
func NewRouter(
	uc port.EcoparqueUseCase,
	loggingMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Healthcheck
	mux.Handle("GET /api/v1/health", parqueHTTP.HealthHandler())

	// CRUD Ecoparques
	mux.Handle("POST /api/v1/parques", parqueHTTP.CreateParqueHandler(uc))
	mux.Handle("GET /api/v1/parques", parqueHTTP.ListParquesHandler(uc))
	mux.Handle("GET /api/v1/parques/{id}", parqueHTTP.GetParqueHandler(uc))
	mux.Handle("PUT /api/v1/parques/{id}", parqueHTTP.UpdateParqueHandler(uc))
	mux.Handle("DELETE /api/v1/parques/{id}", parqueHTTP.DeleteParqueHandler(uc))

	mux.Handle("GET /swagger/", parqueHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", parqueHTTP.SwaggerJSONHandler())

	var handler http.Handler = mux
	if loggingMiddleware != nil {
		handler = loggingMiddleware(handler)
	}

	return handler
}
