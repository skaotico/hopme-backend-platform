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
	Level          LogLevel
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// New crea y configura el logger estructurado JSON del servicio.
func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: strings.ToLower(cfg.Environment) == "development",
	}

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
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithContext inyecta un logger en el contexto.
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithRequestID retorna un nuevo logger hijo con el request_id como atributo fijo.
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}

// WithComponent retorna un nuevo logger hijo con el componente como atributo fijo.
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return logger.With(slog.String("component", component))
}

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

func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, a := range attrs {
		result[i] = a
	}
	return result
}
