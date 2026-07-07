package response

// ==============================================================================
// CÓDIGOS DE ERROR INTERNOS DEL SISTEMA
// ==============================================================================
const (
	// Sistema General (SYS)
	ErrCodeInternal    = "SYS_001" // Error interno del servidor
	ErrCodeInvalidJSON = "SYS_002" // Estructura JSON de petición inválida o malformada
	ErrCodeInvalidID   = "SYS_003" // ID inválido o malformado
)
