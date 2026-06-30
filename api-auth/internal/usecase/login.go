package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/observability"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginUseCase struct {
	userRepo         port.UserRepository
	cacheRepo        port.CacheRepository
	refreshTokenRepo port.RefreshTokenRepository
	jwtSecret        string
	jwtExpiration    time.Duration
	rtExpiration     time.Duration
	maxAttempts      int
	blockMinutes     int
	logger           *slog.Logger
}

// CustomClaims representa la información inyectada en el JWT
type CustomClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// NewLoginUseCase crea una nueva instancia del caso de uso de login
func NewLoginUseCase(
	userRepo port.UserRepository,
	cacheRepo port.CacheRepository,
	refreshTokenRepo port.RefreshTokenRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
	rtExpiration time.Duration,
	maxAttempts int,
	blockMinutes int,
	logger *slog.Logger,
) port.LoginUseCase {
	return &loginUseCase{
		userRepo:         userRepo,
		cacheRepo:        cacheRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
		jwtExpiration:    jwtExpiration,
		rtExpiration:     rtExpiration,
		maxAttempts:      maxAttempts,
		blockMinutes:     blockMinutes,
		logger:           observability.WithComponent(logger, "usecase.login"),
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, email, password, ipAddress, deviceInfo string) (string, string, error) {
	log := observability.FromContext(ctx)
	if log == slog.Default() {
		log = uc.logger
	}

	log.Debug("[LOGIN] Iniciando proceso de autenticación",
		slog.String("email", email),
		slog.String("ip_address", ipAddress),
	)

	// 0. Prevención de fuerza bruta (Rate limiting de Login)
	log.Debug("[LOGIN] Verificando rate limit de intentos fallidos", slog.String("ip_address", ipAddress))
	attempts, err := uc.cacheRepo.IncrLoginAttempt(ctx, ipAddress, uc.blockMinutes*60)
	if err == nil && attempts > uc.maxAttempts {
		log.Info("[LOGIN] Acceso bloqueado por rate limit (demasiados intentos fallidos)",
			slog.String("ip_address", ipAddress),
			slog.Int("attempts", attempts),
			slog.Int("max_attempts", uc.maxAttempts),
		)
		return "", "", model.ErrUnauthorized // Idealmente HTTP 429 Too Many Requests
	}
	log.Debug("[LOGIN] Rate limit OK", slog.String("ip_address", ipAddress), slog.Int("attempts", attempts))

	// 1. Buscar usuario por email
	log.Debug("[LOGIN] Buscando usuario por email", slog.String("email", email))
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		log.Info("[LOGIN] Credenciales inválidas: usuario no encontrado", slog.String("email", email))
		return "", "", model.ErrInvalidCredentials
	}

	// 2. Verificar si el usuario está activo
	if !user.IsActive {
		log.Info("[LOGIN] Intento de acceso con cuenta inactiva",
			slog.String("user_id", user.ID),
			slog.String("email", email),
		)
		return "", "", model.ErrUnauthorized
	}

	// 3. Comparar contraseñas hash
	log.Debug("[LOGIN] Verificando contraseña con bcrypt", slog.String("user_id", user.ID))
	// TODO: ELIMINAR ESTE LOG EN PRODUCCIÓN. Es extremadamente inseguro ya que expone contraseñas en texto plano y hashes.
	log.Warn("[DEBUG-TEMP] Comparando contraseñas",
		slog.String("user_id", user.ID),
		slog.String("stored_hash", user.PasswordHash),
		slog.String("incoming_password", password),
	)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Info("[LOGIN] Credenciales inválidas: contraseña incorrecta",
			slog.String("user_id", user.ID),
			slog.String("email", email),
			slog.String("ip_address", ipAddress),
		)
		return "", "", model.ErrInvalidCredentials
	}

	// 4. Mapear roles a strings
	roleCodes := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleCodes[i] = r.Code
	}
	log.Debug("[LOGIN] Roles del usuario mapeados",
		slog.String("user_id", user.ID),
		slog.Any("roles", roleCodes),
	)

	// 5. Mapear todos los permisos
	permissionCodes := user.GetAllPermissions()
	log.Debug("[LOGIN] Permisos del usuario obtenidos",
		slog.String("user_id", user.ID),
		slog.Int("permission_count", len(permissionCodes)),
	)

	// 6. Generar JTI único para este JWT
	jti := uuid.New().String()
	log.Debug("[LOGIN] JTI generado para el JWT", slog.String("jti", jti))

	// 7. Configurar los claims del token JWT
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

	// 8. Firmar el JWT
	log.Debug("[LOGIN] Firmando el JWT con HS256", slog.String("user_id", user.ID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		log.Error("[LOGIN] Error al firmar el JWT",
			slog.String("user_id", user.ID),
			slog.Any("error", err),
		)
		return "", "", err
	}

	// 9. Generar Refresh Token Opaco
	log.Debug("[LOGIN] Generando refresh token opaco...", slog.String("user_id", user.ID))
	plainRefreshToken := uuid.New().String() + "-" + uuid.New().String()

	// Hashear el Refresh Token antes de guardar usando SHA-256
	hasher := sha256.New()
	hasher.Write([]byte(plainRefreshToken))
	rtHash := hex.EncodeToString(hasher.Sum(nil))

	rt := &model.RefreshToken{
		UserID:     user.ID,
		TokenHash:  rtHash,
		ExpiresAt:  time.Now().Add(uc.rtExpiration),
		Revoked:    false,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	}

	// Guardar en DB
	log.Debug("[LOGIN] Persistiendo refresh token en base de datos...", slog.String("user_id", user.ID))
	err = uc.refreshTokenRepo.Create(ctx, rt)
	if err != nil {
		log.Error("[LOGIN] Error al persistir el refresh token",
			slog.String("user_id", user.ID),
			slog.Any("error", err),
		)
		return "", "", err
	}

	log.Info("[LOGIN] Autenticación exitosa",
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("email", email),
		slog.String("ip_address", ipAddress),
		slog.String("jti", jti),
	)

	return tokenStr, plainRefreshToken, nil
}
