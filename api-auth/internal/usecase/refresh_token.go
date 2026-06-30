package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/observability"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type refreshTokenUseCase struct {
	refreshTokenRepo port.RefreshTokenRepository
	userRepo         port.UserRepository
	jwtSecret        string
	jwtExpiration    time.Duration
	rtExpiration     time.Duration
	logger           *slog.Logger
}

// NewRefreshTokenUseCase crea el caso de uso para refrescar tokens
func NewRefreshTokenUseCase(
	refreshTokenRepo port.RefreshTokenRepository,
	userRepo port.UserRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
	rtExpiration time.Duration,
	logger *slog.Logger,
) port.RefreshTokenUseCase {
	return &refreshTokenUseCase{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtExpiration:    jwtExpiration,
		rtExpiration:     rtExpiration,
		logger:           observability.WithComponent(logger, "usecase.refresh_token"),
	}
}

func (uc *refreshTokenUseCase) Execute(ctx context.Context, refreshToken, ipAddress, deviceInfo string) (string, string, error) {
	log := observability.FromContext(ctx)
	if log == slog.Default() {
		log = uc.logger
	}

	log.Debug("[REFRESH] Iniciando rotación de refresh token",
		slog.String("ip_address", ipAddress),
	)

	// 1. Hashear el token recibido para buscarlo
	log.Debug("[REFRESH] Hasheando token recibido para búsqueda en BD...")
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	rtHash := hex.EncodeToString(hasher.Sum(nil))

	// 2. Buscar en la base de datos
	log.Debug("[REFRESH] Buscando refresh token en la base de datos...")
	rt, err := uc.refreshTokenRepo.FindByHash(ctx, rtHash)
	if err != nil || rt == nil {
		log.Info("[REFRESH] Refresh token no encontrado en la base de datos",
			slog.String("ip_address", ipAddress),
			slog.Any("error", err),
		)
		return "", "", errors.New("invalid refresh token")
	}

	log.Debug("[REFRESH] Token encontrado",
		slog.String("token_id", rt.ID),
		slog.String("user_id", rt.UserID),
		slog.Bool("revoked", rt.Revoked),
	)

	// 3. Verificar si el token fue revocado (alerta de seguridad / rotación)
	if rt.Revoked {
		// Posible robo de token. Revocamos toda la familia de tokens del usuario.
		log.Info("[REFRESH] ⚠️  ALERTA DE SEGURIDAD: token ya revocado fue reutilizado, revocando familia completa",
			slog.String("user_id", rt.UserID),
			slog.String("token_id", rt.ID),
			slog.String("ip_address", ipAddress),
		)
		_ = uc.refreshTokenRepo.RevokeFamily(ctx, rt.UserID)
		return "", "", errors.New("token reused - security alert")
	}

	// 4. Verificar expiración
	if !rt.IsValid() {
		log.Info("[REFRESH] Refresh token expirado",
			slog.String("user_id", rt.UserID),
			slog.String("token_id", rt.ID),
			slog.Time("expired_at", rt.ExpiresAt),
		)
		return "", "", errors.New("refresh token expired")
	}

	// 5. Inmediatamente rotar el token (revocar el actual)
	log.Debug("[REFRESH] Revocando token actual para rotación",
		slog.String("token_id", rt.ID),
		slog.String("user_id", rt.UserID),
	)
	_ = uc.refreshTokenRepo.Revoke(ctx, rt.ID)

	// 6. Obtener el usuario
	log.Debug("[REFRESH] Cargando datos actualizados del usuario", slog.String("user_id", rt.UserID))
	user, err := uc.userRepo.FindByID(ctx, rt.UserID)
	if err != nil || user == nil || !user.IsActive {
		log.Info("[REFRESH] Usuario no encontrado o inactivo al refrescar token",
			slog.String("user_id", rt.UserID),
			slog.Any("error", err),
		)
		return "", "", errors.New("user not found or inactive")
	}

	// 7. Generar nuevos tokens (Access y Refresh)
	roleCodes := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleCodes[i] = r.Code
	}

	permissionCodes := user.GetAllPermissions()
	jti := uuid.New().String()

	log.Debug("[REFRESH] Generando nuevo JWT",
		slog.String("user_id", user.ID),
		slog.String("jti", jti),
		slog.Any("roles", roleCodes),
	)

	claims := CustomClaims{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roleCodes,
		Permissions: permissionCodes,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(uc.jwtExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "homelab-auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		log.Error("[REFRESH] Error al firmar el nuevo JWT",
			slog.String("user_id", user.ID),
			slog.Any("error", err),
		)
		return "", "", err
	}

	// 8. Generar nuevo Refresh Token
	log.Debug("[REFRESH] Generando nuevo refresh token opaco...", slog.String("user_id", user.ID))
	newPlainRefreshToken := uuid.New().String() + "-" + uuid.New().String()
	newHasher := sha256.New()
	newHasher.Write([]byte(newPlainRefreshToken))
	newRtHash := hex.EncodeToString(newHasher.Sum(nil))

	newRt := &model.RefreshToken{
		UserID:     user.ID,
		TokenHash:  newRtHash,
		ExpiresAt:  time.Now().Add(uc.rtExpiration),
		Revoked:    false,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	}

	log.Debug("[REFRESH] Persistiendo nuevo refresh token en base de datos...", slog.String("user_id", user.ID))
	err = uc.refreshTokenRepo.Create(ctx, newRt)
	if err != nil {
		log.Error("[REFRESH] Error al persistir el nuevo refresh token",
			slog.String("user_id", user.ID),
			slog.Any("error", err),
		)
		return "", "", err
	}

	log.Info("[REFRESH] Rotación de tokens completada exitosamente",
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("new_jti", jti),
		slog.String("ip_address", ipAddress),
	)

	return tokenStr, newPlainRefreshToken, nil
}
