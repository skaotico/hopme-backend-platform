package port

import (
	"context"
	"c4-auth/internal/domain/model"
)

// RegisterUseCase define el Puerto de Entrada para el registro de nuevos usuarios
type RegisterUseCase interface {
	Execute(ctx context.Context, username, email, password string) (*model.User, error)
}

// LoginUseCase define el Puerto de Entrada para el inicio de sesión y emisión de tokens
type LoginUseCase interface {
	Execute(ctx context.Context, email, password, ipAddress, deviceInfo string) (string, string, error) // Retorna accessToken, refreshToken
}

// RefreshTokenUseCase define el Puerto de Entrada para la rotación de tokens
type RefreshTokenUseCase interface {
	Execute(ctx context.Context, refreshToken, ipAddress, deviceInfo string) (string, string, error) // Retorna accessToken, refreshToken
}

// LogoutUseCase define el Puerto de Entrada para cerrar sesión
type LogoutUseCase interface {
	Execute(ctx context.Context, jti string, timeToLive int, refreshToken string) error
}
