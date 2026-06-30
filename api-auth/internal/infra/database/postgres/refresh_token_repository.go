package postgres

import (
	"context"
	"database/sql"
	"errors"

	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
)

type refreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository crea una instancia del repositorio para tokens de refresco
func NewRefreshTokenRepository(db *sql.DB) port.RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	query := `
		INSERT INTO auth.refresh_tokens (user_id, token_hash, expires_at, revoked, device_info, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt, token.Revoked, token.DeviceInfo, token.IPAddress,
	).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)

	return err
}

func (r *refreshTokenRepository) FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked, device_info, ip_address, created_at, updated_at
		FROM auth.refresh_tokens
		WHERE token_hash = $1
	`
	token := &model.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.Revoked,
		&token.DeviceInfo, &token.IPAddress, &token.CreatedAt, &token.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("refresh token not found") // Idealmente usar model.ErrRefreshTokenNotFound si existiera
		}
		return nil, err
	}

	return token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id string) error {
	query := `
		UPDATE auth.refresh_tokens
		SET revoked = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *refreshTokenRepository) RevokeFamily(ctx context.Context, userID string) error {
	query := `
		UPDATE auth.refresh_tokens
		SET revoked = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND revoked = FALSE
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
