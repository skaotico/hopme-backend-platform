// Package observability centraliza la configuración del sistema de logs estructurados
// del servicio Catalogo. Usa log/slog como implementación nativa.
package observability

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const loggerKey contextKey = "logger"

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type Config struct {
	Level          LogLevel
	ServiceName    string
	ServiceVersion string
	Environment    string
}

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

	slog.SetDefault(logger)

	return logger
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}

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
