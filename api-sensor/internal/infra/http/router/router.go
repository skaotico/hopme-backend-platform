package router

import (
	"log/slog"
	"net/http"

	sensorHTTP "c4-sensor/internal/infra/http"
	"c4-sensor/internal/domain/port"
)

// NewRouter crea y configura un multiplexor de peticiones para API Sensor.
func NewRouter(
	sensorUC port.SensorUseCase,
	lecturaUC port.LecturaUseCase,
	alertaUC port.AlertaUseCase,
	log *slog.Logger,
	loggingMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Healthcheck
	mux.Handle("GET /api/v1/health", sensorHTTP.HealthHandler())

	// CRUD Sensores
	mux.Handle("POST /api/v1/sensores", sensorHTTP.CreateSensorHandler(sensorUC, log))
	mux.Handle("GET /api/v1/sensores", sensorHTTP.ListSensoresHandler(sensorUC, log))
	mux.Handle("GET /api/v1/sensores/{id}", sensorHTTP.GetSensorHandler(sensorUC, log))
	mux.Handle("PUT /api/v1/sensores/{id}", sensorHTTP.UpdateSensorHandler(sensorUC, log))
	mux.Handle("DELETE /api/v1/sensores/{id}", sensorHTTP.DeleteSensorHandler(sensorUC, log))

	// Lecturas del Sensor
	mux.Handle("POST /api/v1/sensores/{id}/lecturas", sensorHTTP.CreateLecturaHandler(lecturaUC, log))
	mux.Handle("GET /api/v1/sensores/{id}/lecturas", sensorHTTP.ListLecturasHandler(lecturaUC, log))

	// Alertas de Sensores
	mux.Handle("POST /api/v1/sensores/{id}/alertas", sensorHTTP.CreateAlertaHandler(alertaUC, log))
	mux.Handle("GET /api/v1/alertas", sensorHTTP.ListAlertasHandler(alertaUC, log))
	mux.Handle("PUT /api/v1/alertas/{id}/estado", sensorHTTP.UpdateAlertaStatusHandler(alertaUC, log))

	// Swagger UI
	mux.Handle("GET /swagger/", sensorHTTP.SwaggerHandler())
	mux.Handle("GET /swagger/doc.json", sensorHTTP.SwaggerJSONHandler())

	var handler http.Handler = mux
	if loggingMiddleware != nil {
		handler = loggingMiddleware(handler)
	}

	return handler
}
