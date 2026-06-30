package port

import (
	"context"
	"c4-auth/internal/domain/model"
)

// UserRepository define el Puerto de Salida para la persistencia de usuarios, roles y permisos
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	AssignRole(ctx context.Context, userID, roleCode string) error
	FindRoleByCode(ctx context.Context, code string) (*model.Role, error)
}

// RefreshTokenRepository define el Puerto de Salida para los tokens de refresco
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeFamily(ctx context.Context, userID string) error
}

// CacheRepository define el Puerto de Salida para la caché y blacklist
type CacheRepository interface {
	BlacklistJWT(ctx context.Context, jti string, ttlSeconds int) error
	IsJWTBlacklisted(ctx context.Context, jti string) (bool, error)
	IncrLoginAttempt(ctx context.Context, ip string, ttlSeconds int) (int, error)
	GetProfile(ctx context.Context, userID string) (string, error)
	SetProfile(ctx context.Context, userID string, profileJSON string, ttlSeconds int) error
	DeleteProfile(ctx context.Context, userID string) error
}
