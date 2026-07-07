package router

import (
	"net/http"

	arbolHTTP "c4-arbol/internal/infra/http"
	"c4-arbol/internal/domain/port"
)

// NewRouter crea y configura un multiplexor de peticiones para API Arbol.
func NewRouter(
	arbolUC port.ArbolUseCase,
	historialUC port.HistorialUseCase,
	medicionUC port.MedicionUseCase,
	loggingMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Healthcheck
	mux.Handle("GET /api/v1/health", arbolHTTP.HealthHandler())

	// CRUD Árboles
	mux.Handle("POST /api/v1/arboles", arbolHTTP.CreateArbolHandler(arbolUC))
	mux.Handle("GET /api/v1/arboles", arbolHTTP.ListArbolesHandler(arbolUC))
	mux.Handle("GET /api/v1/arboles/{id}", arbolHTTP.GetArbolHandler(arbolUC))
	mux.Handle("PUT /api/v1/arboles/{id}", arbolHTTP.UpdateArbolHandler(arbolUC))
	mux.Handle("DELETE /api/v1/arboles/{id}", arbolHTTP.DeleteArbolHandler(arbolUC))

	// Historial de Estado
	mux.Handle("POST /api/v1/arboles/{id}/historial", arbolHTTP.CreateHistorialHandler(historialUC))
	mux.Handle("GET /api/v1/arboles/{id}/historial", arbolHTTP.ListHistorialHandler(historialUC))

	// Mediciones Dasométricas
	mux.Handle("POST /api/v1/arboles/{id}/mediciones", arbolHTTP.CreateMedicionHandler(medicionUC))
	mux.Handle("GET /api/v1/arboles/{id}/mediciones", arbolHTTP.ListMedicionesHandler(medicionUC))

	// Swagger UI
	mux.Handle("GET /swagger/", arbolHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", arbolHTTP.SwaggerJSONHandler())

	var handler http.Handler = mux
	if loggingMiddleware != nil {
		handler = loggingMiddleware(handler)
	}

	return handler
}
