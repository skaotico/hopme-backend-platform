package response

// ==============================================================================
// CÓDIGOS DE ERROR INTERNOS DEL SISTEMA (MAREADOS POR API)
// ==============================================================================
const (
	// Sistema General (SYS)
	ErrCodeInternal    = "SYS_001" // Error interno del servidor
	ErrCodeInvalidJSON = "SYS_002" // Estructura JSON de petición inválida o malformada

	// Autenticación & Autorización (AUTH)
	ErrCodeInvalidEmail       = "AUTH_001" // Correo electrónico inválido
	ErrCodeInvalidUsername    = "AUTH_002" // Nombre de usuario inválido
	ErrCodeInvalidPassword    = "AUTH_003" // Contraseña no cumple longitud mínima
	ErrCodeUserExists         = "AUTH_004" // Nombre de usuario o correo ya registrado
	ErrCodeInvalidCredentials = "AUTH_005" // Correo o contraseña incorrectos
	ErrCodeUserDisabled       = "AUTH_006" // Cuenta desactivada
	ErrCodeUnauthorized       = "AUTH_007" // Token JWT ausente, inválido o expirado
)
