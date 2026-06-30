package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"

	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/observability"
)

type logoutUseCase struct {
	cacheRepo        port.CacheRepository
	refreshTokenRepo port.RefreshTokenRepository
	logger           *slog.Logger
}

// NewLogoutUseCase crea el caso de uso de logout
func NewLogoutUseCase(
	cacheRepo port.CacheRepository,
	refreshTokenRepo port.RefreshTokenRepository,
	logger *slog.Logger,
) port.LogoutUseCase {
	return &logoutUseCase{
		cacheRepo:        cacheRepo,
		refreshTokenRepo: refreshTokenRepo,
		logger:           observability.WithComponent(logger, "usecase.logout"),
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, jti string, timeToLive int, refreshToken string) error {
	log := observability.FromContext(ctx)
	if log == slog.Default() {
		log = uc.logger
	}

	log.Debug("[LOGOUT] Iniciando proceso de logout",
		slog.String("jti", jti),
		slog.Int("ttl_seconds", timeToLive),
		slog.Bool("has_refresh_token", refreshToken != ""),
	)

	// 1. Añadir a la Blacklist el JWT actual (si queda tiempo de vida)
	if timeToLive > 0 && jti != "" {
		log.Debug("[LOGOUT] Añadiendo JWT a la blacklist en cache",
			slog.String("jti", jti),
			slog.Int("ttl_seconds", timeToLive),
		)
		if err := uc.cacheRepo.BlacklistJWT(ctx, jti, timeToLive); err != nil {
			log.Info("[LOGOUT] No se pudo añadir JWT a la blacklist (no crítico)",
				slog.String("jti", jti),
				slog.Any("error", err),
			)
		} else {
			log.Debug("[LOGOUT] JWT añadido a la blacklist correctamente", slog.String("jti", jti))
		}
	} else {
		log.Debug("[LOGOUT] JWT omitido de la blacklist (TTL <= 0 o JTI vacío)",
			slog.String("jti", jti),
			slog.Int("ttl_seconds", timeToLive),
		)
	}

	// 2. Si se proporcionó un Refresh Token, revocarlo en la base de datos
	if refreshToken != "" {
		log.Debug("[LOGOUT] Hasheando refresh token para revocación en BD...")
		hasher := sha256.New()
		hasher.Write([]byte(refreshToken))
		rtHash := hex.EncodeToString(hasher.Sum(nil))

		rt, err := uc.refreshTokenRepo.FindByHash(ctx, rtHash)
		if err != nil {
			log.Info("[LOGOUT] Refresh token no encontrado al intentar revocar (puede ya estar revocado)",
				slog.Any("error", err),
			)
		} else if rt != nil {
			log.Debug("[LOGOUT] Revocando refresh token en la base de datos",
				slog.String("token_id", rt.ID),
				slog.String("user_id", rt.UserID),
			)
			if err := uc.refreshTokenRepo.Revoke(ctx, rt.ID); err != nil {
				log.Error("[LOGOUT] Error al revocar refresh token",
					slog.String("token_id", rt.ID),
					slog.Any("error", err),
				)
			} else {
				log.Debug("[LOGOUT] Refresh token revocado correctamente",
					slog.String("token_id", rt.ID),
					slog.String("user_id", rt.UserID),
				)
			}
		}
	} else {
		log.Debug("[LOGOUT] No se proporcionó refresh token, omitiendo revocación en BD")
	}

	log.Info("[LOGOUT] Logout completado",
		slog.String("jti", jti),
		slog.Bool("jwt_blacklisted", timeToLive > 0 && jti != ""),
		slog.Bool("refresh_token_revoked", refreshToken != ""),
	)

	return nil
}
