package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// Timeout crea un middleware que limita el tiempo de ejecución de una petición
// utilizando el http.TimeoutHandler nativo de Go.
func Timeout(duration time.Duration, routeName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Estructura de error JSON consistente con la API
		jsonMsg := fmt.Sprintf(
			`{"success":false,"error":{"code":"SYS_003","message":"El tiempo de espera de la petición para %s ha expirado","details":"Timeout de %s alcanzado"}}`,
			routeName,
			duration.String(),
		)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Establecemos de antemano el tipo de contenido como JSON para que si se produce
			// un timeout, el cuerpo del mensaje de error se interprete correctamente como tal.
			w.Header().Set("Content-Type", "application/json")
			
			// Creamos el TimeoutHandler nativo
			handler := http.TimeoutHandler(next, duration, jsonMsg)
			
			// Servimos la petición
			handler.ServeHTTP(w, r)
		})
	}
}
