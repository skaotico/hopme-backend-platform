// Package observability centraliza la configuración del sistema de logs estructurados
// del servicio Auth. Usa log/slog (stdlib Go 1.21+) como implementación nativa,
// emitiendo JSON listo para ser consumido por Loki, Grafana, OpenTelemetry, etc.
package observability

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// contextKey es un tipo privado para evitar colisiones de claves en el contexto.
type contextKey string

const loggerKey contextKey = "logger"

// LogLevel define los niveles de log soportados por el servicio.
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config define las opciones de configuración del logger.
type Config struct {
	// Level es el nivel mínimo de log a emitir ("debug", "info", "warn", "error").
	// Por defecto: "info"
	Level LogLevel

	// ServiceName se inyecta como atributo fijo en todos los logs para identificar
	// el microservicio en sistemas de observabilidad centralizados (Loki, Grafana).
	ServiceName string

	// ServiceVersion se inyecta como atributo fijo para correlacionar logs con
	// versiones de despliegue.
	ServiceVersion string

	// Environment indica el entorno ("development", "staging", "production").
	// En "development" los logs incluyen el caller (archivo:línea) para depuración.
	Environment string
}

// New crea y configura el logger estructurado JSON del servicio.
// Registra el logger como default global de slog para compatibilidad
// con código legacy que use slog.Info/Error directamente.
func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level: level,
		// AddSource = true en desarrollo para ver archivo:línea en cada log
		AddSource: strings.ToLower(cfg.Environment) == "development",
	}

	// Atributos fijos que se incluyen en TODOS los logs del servicio.
	// Fundamentales para filtrar en Loki/Grafana por servicio y versión.
	baseAttrs := []slog.Attr{
		slog.String("service", cfg.ServiceName),
		slog.String("version", cfg.ServiceVersion),
		slog.String("env", cfg.Environment),
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler).With(attrsToAny(baseAttrs)...)

	// Registrar como logger global para compatibilidad con slog.Info/Error
	slog.SetDefault(logger)

	return logger
}

// FromContext recupera el logger inyectado en el contexto.
// Si no existe, retorna el logger global (siempre válido).
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithContext inyecta un logger en el contexto para propagación
// a lo largo de una cadena de llamadas (handler → usecase → repository).
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithRequestID retorna un nuevo logger hijo con el request_id como atributo fijo.
// Úsalo al inicio de cada handler HTTP para correlacionar todos los logs de esa request.
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}

// WithUserID retorna un nuevo logger hijo con el user_id como atributo fijo.
// Úsalo al autenticar al usuario para correlacionar logs de sesión.
func WithUserID(logger *slog.Logger, userID string) *slog.Logger {
	return logger.With(slog.String("user_id", userID))
}

// WithComponent retorna un nuevo logger hijo con el componente como atributo fijo.
// Permite filtrar logs por capa: "usecase", "repository", "handler", "middleware".
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return logger.With(slog.String("component", component))
}

// parseLevel convierte el string de configuración al nivel slog correspondiente.
func parseLevel(l LogLevel) slog.Level {
	switch strings.ToLower(string(l)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// attrsToAny convierte []slog.Attr a []any para uso en .With(...)
func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, a := range attrs {
		result[i] = a
	}
	return result
}
