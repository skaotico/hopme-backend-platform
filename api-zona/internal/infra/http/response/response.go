package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse representa el envoltorio único y genérico de respuesta para el Front
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError representa la información detallada del error interno
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// JSON escribe una respuesta JSON estructurada directa
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Success responde una solicitud exitosa utilizando el envoltorio estándar
func Success(w http.ResponseWriter, status int, data interface{}) {
	res := APIResponse{
		Success: true,
		Data:    data,
	}
	JSON(w, status, res)
}

// Failure responde una solicitud fallida utilizando el envoltorio estándar y código interno
func Failure(w http.ResponseWriter, status int, errorCode string, message string, details string) {
	res := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
	JSON(w, status, res)
}
