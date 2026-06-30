package usecase

import (
	"context"
	"log/slog"

	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
	"c4-auth/internal/infra/observability"

	"golang.org/x/crypto/bcrypt"
)

type registerUseCase struct {
	userRepo port.UserRepository
	logger   *slog.Logger
}

// NewRegisterUseCase crea una nueva instancia del caso de uso de registro
func NewRegisterUseCase(userRepo port.UserRepository, logger *slog.Logger) port.RegisterUseCase {
	return &registerUseCase{
		userRepo: userRepo,
		logger:   observability.WithComponent(logger, "usecase.register"),
	}
}

func (uc *registerUseCase) Execute(ctx context.Context, username, email, password string) (*model.User, error) {
	log := observability.FromContext(ctx)
	if log == slog.Default() {
		log = uc.logger
	}

	log.Debug("[REGISTER] Iniciando proceso de registro",
		slog.String("username", username),
		slog.String("email", email),
	)

	// 1. Crear la entidad de dominio con los nuevos campos de la BD V1
	user := &model.User{
		Username:    username,
		Email:       email,
		DisplayName: username, // Por defecto el DisplayName coincide con el username
		IsActive:    true,
		IsVerified:  false, // Por defecto el usuario no está verificado al crearse
	}

	// 2. Validaciones puras del dominio
	if err := user.Validate(); err != nil {
		log.Debug("[REGISTER] Validación de dominio fallida",
			slog.String("username", username),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return nil, err
	}

	if len(password) < 6 {
		log.Debug("[REGISTER] Contraseña inválida: longitud insuficiente",
			slog.String("username", username),
		)
		return nil, model.ErrInvalidPassword
	}

	log.Debug("[REGISTER] Verificando si el email ya existe", slog.String("email", email))

	// 3. Verificar si el usuario ya existe en base a email o username
	existingEmail, _ := uc.userRepo.FindByEmail(ctx, user.Email)
	if existingEmail != nil {
		log.Info("[REGISTER] Intento de registro con email ya existente",
			slog.String("email", email),
		)
		return nil, model.ErrUserAlreadyExists
	}

	log.Debug("[REGISTER] Verificando si el username ya existe", slog.String("username", username))

	existingUser, _ := uc.userRepo.FindByUsername(ctx, user.Username)
	if existingUser != nil {
		log.Info("[REGISTER] Intento de registro con username ya existente",
			slog.String("username", username),
		)
		return nil, model.ErrUserAlreadyExists
	}

	// 4. Hashing de contraseña (bcrypt)
	log.Debug("[REGISTER] Generando hash bcrypt de la contraseña...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("[REGISTER] Error al generar hash de contraseña",
			slog.String("username", username),
			slog.Any("error", err),
		)
		return nil, err
	}
	user.PasswordHash = string(hashedPassword)

	// 5. Guardar en base de datos (esquema auth)
	log.Debug("[REGISTER] Persistiendo usuario en base de datos...", slog.String("username", username))
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		log.Error("[REGISTER] Error al crear usuario en la base de datos",
			slog.String("username", username),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return nil, err
	}

	log.Debug("[REGISTER] Usuario creado correctamente, asignando rol predeterminado...",
		slog.String("user_id", user.ID),
		slog.String("username", username),
	)

	// 6. Asignar rol predeterminado 'USER' (esquema iam)
	defaultRole, err := uc.userRepo.FindRoleByCode(ctx, "USER")
	if err == nil && defaultRole != nil {
		_ = uc.userRepo.AssignRole(ctx, user.ID, defaultRole.Code)
		user.Roles = append(user.Roles, *defaultRole)
		log.Debug("[REGISTER] Rol predeterminado asignado",
			slog.String("user_id", user.ID),
			slog.String("role", defaultRole.Code),
		)
	} else {
		log.Info("[REGISTER] No se encontró rol predeterminado 'USER', usuario creado sin rol",
			slog.String("user_id", user.ID),
		)
	}

	log.Info("[REGISTER] Registro completado exitosamente",
		slog.String("user_id", user.ID),
		slog.String("username", username),
		slog.String("email", email),
	)

	return user, nil
}
