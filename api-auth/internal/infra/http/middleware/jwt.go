package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/observability"
	"c4-auth/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

// JWTAuth es un middleware nativo que intercepta solicitudes, valida el JWT
// y verifica que el JTI no esté en la blacklist de Redis (logout previo).
func JWTAuth(jwtSecret string, cacheRepo port.CacheRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := observability.FromContext(r.Context())

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Debug("[JWT] Authorization header ausente")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":{"code":"AUTH_007","message":"encabezado Authorization requerido"}}`))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Debug("[JWT] Formato de Authorization header inválido")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":{"code":"AUTH_007","message":"formato de token inválido, debe ser Bearer <token>"}}`))
				return
			}

			tokenString := parts[1]
			claims := &usecase.CustomClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				log.Debug("[JWT] Token inválido o expirado", slog.Any("error", err))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":{"code":"AUTH_007","message":"token de acceso inválido, expirado o alterado"}}`))
				return
			}

			// Verificar si el JTI del JWT está en la blacklist de Redis (logout previo)
			if claims.ID != "" {
				log.Debug("[JWT] Verificando JTI en blacklist", slog.String("jti", claims.ID))

				blacklisted, err := cacheRepo.IsJWTBlacklisted(r.Context(), claims.ID)
				if err != nil {
					// Fallo en Redis: por seguridad, loggear pero dejar pasar (fail-open).
					// En entornos críticos cambiar a fail-closed (rechazar).
					log.Info("[JWT] No se pudo verificar blacklist en Redis, permitiendo acceso (fail-open)",
						slog.String("jti", claims.ID),
						slog.Any("error", err),
					)
				} else if blacklisted {
					log.Info("[JWT] Token en blacklist rechazado",
						slog.String("jti", claims.ID),
						slog.String("user_id", claims.UserID),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(`{"success":false,"error":{"code":"AUTH_008","message":"sesión cerrada, por favor inicia sesión nuevamente"}}`))
					return
				}
			}

			log.Debug("[JWT] Token válido, inyectando claims en contexto",
				slog.String("jti", claims.ID),
				slog.String("user_id", claims.UserID),
			)

			// Inyectar los claims decodificados en el contexto de la petición
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims recupera de forma segura los claims del usuario desde el contexto
func GetUserClaims(ctx context.Context) (*usecase.CustomClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*usecase.CustomClaims)
	return claims, ok
}
