package response

// RegisterRequest representa la solicitud HTTP para registrar un usuario
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest representa la solicitud HTTP para iniciar sesión
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse representa el DTO de respuesta con un token JWT y Refresh Token
type TokenResponse struct {
	Token string `json:"token"`
}

// RefreshRequest representa la solicitud HTTP para rotar el token
type RefreshRequest struct {
	// El refresh_token viaja ahora en Cookie, por lo que el body puede estar vacío
}

// LogoutRequest representa la solicitud HTTP para cerrar sesión
type LogoutRequest struct {
	// El refresh_token viaja ahora en Cookie, por lo que el body puede estar vacío
}
