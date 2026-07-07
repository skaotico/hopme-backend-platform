package router

import (
	"net/http"

	catalogoHTTP "c4-catalogo/internal/infra/http"
)

// NewRouter crea y configura un multiplexor de peticiones nativo para toda la API de catálogo
func NewRouter(
	especieArbolH *catalogoHTTP.EspecieArbolHandler,
	estadoArbolH *catalogoHTTP.EstadoArbolHandler,
	estadoEstanqueH *catalogoHTTP.EstadoEstanqueHandler,
	estadoAguaH *catalogoHTTP.EstadoAguaHandler,
	tipoSensorH *catalogoHTTP.TipoSensorHandler,
	loggingMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Swagger Docs
	mux.Handle("GET /swagger/", catalogoHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", catalogoHTTP.SwaggerJSONHandler())

	// Endpoint público de Healthcheck
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ok"}}`))
	})

	// EspecieArbol
	mux.Handle("GET /api/v1/catalogo/especies", especieArbolH.List())
	mux.Handle("POST /api/v1/catalogo/especies", especieArbolH.Create())
	mux.Handle("GET /api/v1/catalogo/especies/{id}", especieArbolH.GetByID())
	mux.Handle("PUT /api/v1/catalogo/especies/{id}", especieArbolH.Update())
	mux.Handle("DELETE /api/v1/catalogo/especies/{id}", especieArbolH.Delete())

	// EstadoArbol
	mux.Handle("GET /api/v1/catalogo/estados-arbol", estadoArbolH.List())
	mux.Handle("POST /api/v1/catalogo/estados-arbol", estadoArbolH.Create())
	mux.Handle("GET /api/v1/catalogo/estados-arbol/{id}", estadoArbolH.GetByID())
	mux.Handle("PUT /api/v1/catalogo/estados-arbol/{id}", estadoArbolH.Update())
	mux.Handle("DELETE /api/v1/catalogo/estados-arbol/{id}", estadoArbolH.Delete())

	// EstadoEstanque
	mux.Handle("GET /api/v1/catalogo/estados-estanque", estadoEstanqueH.List())
	mux.Handle("POST /api/v1/catalogo/estados-estanque", estadoEstanqueH.Create())
	mux.Handle("GET /api/v1/catalogo/estados-estanque/{id}", estadoEstanqueH.GetByID())
	mux.Handle("PUT /api/v1/catalogo/estados-estanque/{id}", estadoEstanqueH.Update())
	mux.Handle("DELETE /api/v1/catalogo/estados-estanque/{id}", estadoEstanqueH.Delete())

	// EstadoAgua
	mux.Handle("GET /api/v1/catalogo/estados-agua", estadoAguaH.List())
	mux.Handle("POST /api/v1/catalogo/estados-agua", estadoAguaH.Create())
	mux.Handle("GET /api/v1/catalogo/estados-agua/{id}", estadoAguaH.GetByID())
	mux.Handle("PUT /api/v1/catalogo/estados-agua/{id}", estadoAguaH.Update())
	mux.Handle("DELETE /api/v1/catalogo/estados-agua/{id}", estadoAguaH.Delete())

	// TipoSensor
	mux.Handle("GET /api/v1/catalogo/tipos-sensor", tipoSensorH.List())
	mux.Handle("POST /api/v1/catalogo/tipos-sensor", tipoSensorH.Create())
	mux.Handle("GET /api/v1/catalogo/tipos-sensor/{id}", tipoSensorH.GetByID())
	mux.Handle("PUT /api/v1/catalogo/tipos-sensor/{id}", tipoSensorH.Update())
	mux.Handle("DELETE /api/v1/catalogo/tipos-sensor/{id}", tipoSensorH.Delete())

	var handler http.Handler = mux
	handler = loggingMiddleware(handler)

	return handler
}
