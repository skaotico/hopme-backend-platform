package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"c4-catalogo/internal/infra/observability"
)

// responseWriterDelegator wraps a standard ResponseWriter to record status codes
type responseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterDelegator) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterDelegator) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Logging middleware logs each request and duration
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate raw Request ID
		bytes := make([]byte, 8)
		_, _ = rand.Read(bytes)
		requestID := hex.EncodeToString(bytes)

		// Set logger with Request ID in context
		logger := observability.WithRequestID(slog.Default(), requestID)
		ctx := observability.WithContext(r.Context(), logger)

		delegator := &responseWriterDelegator{ResponseWriter: w}

		logger.Info("[HTTP] Iniciando petición",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(delegator, r.WithContext(ctx))

		logger.Info("[HTTP] Petición finalizada",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", delegator.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
