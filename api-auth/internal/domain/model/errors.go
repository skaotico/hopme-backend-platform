package model

import "fmt"

// AppError es nuestra estructura de error estándar para la API
type AppError struct {
	Code    string
	Message string
}

// Error implementa la interfaz de error nativa de Go
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Errores de la capa de Auth
var (
	// Errores de Validación (1000)
	ErrUserAlreadyExists = &AppError{Code: "AUTH-1001", Message: "el usuario o correo electrónico ya se encuentra registrado"}

	// Errores de Autenticación / Autorización (2000)
	ErrInvalidCredentials = &AppError{Code: "AUTH-2001", Message: "credenciales inválidas, verifique su correo o contraseña"}
	ErrUnauthorized       = &AppError{Code: "AUTH-2002", Message: "no autorizado para realizar esta acción"}

	// Errores de Negocio / Dominio (3000)
	ErrUserNotFound = &AppError{Code: "AUTH-3001", Message: "el usuario no fue encontrado"}
	ErrRoleNotFound = &AppError{Code: "AUTH-3002", Message: "el rol especificado no existe"}
)
