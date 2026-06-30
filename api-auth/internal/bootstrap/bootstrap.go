package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"c4-auth/internal/infra/cache"
	"c4-auth/internal/infra/config"
	"c4-auth/internal/infra/database"
	"c4-auth/internal/infra/database/postgres"
	"c4-auth/internal/infra/http/middleware"
	"c4-auth/internal/infra/http/router"
	"c4-auth/internal/infra/observability"
	"c4-auth/internal/usecase"

	"github.com/redis/go-redis/v9"
)

// Run centraliza toda la inicialización, inyección de dependencias y arranque del servicio Auth
func Run() {
	// 1. Cargar configuración PRIMERO (necesitamos LogLevel antes de loggear)
	cfg := config.Load()

	// 2. Inicializar el logger estructurado con el paquete de observabilidad
	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "auth-service",
		ServiceVersion: "v1.3",
		Environment:    cfg.Environment,
	})

	logger.Info("[AUTH-SERVICE] Iniciando servicio de autenticación...",
		slog.String("log_level", cfg.LogLevel),
		slog.String("environment", cfg.Environment),
		slog.String("port", cfg.Port),
	)

	// 3. Establecer conexión con PostgreSQL (Pool de SQL nativo)
	logger.Debug("[DATABASE] Construyendo configuración de conexión a PostgreSQL",
		slog.String("host", cfg.DBHost),
		slog.Int("port", cfg.DBPort),
		slog.String("dbname", cfg.DBName),
		slog.String("sslmode", cfg.DBSSLMode),
	)

	dbCfg := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	db, err := database.NewConnection(dbCfg)
	if err != nil {
		logger.Error("[AUTH-SERVICE] Error crítico al conectar a la base de datos",
			slog.Any("error", err),
			slog.String("host", cfg.DBHost),
			slog.Int("port", cfg.DBPort),
		)
		os.Exit(1)
	}
	logger.Info("[DATABASE] Conexión a PostgreSQL establecida correctamente",
		slog.String("host", cfg.DBHost),
		slog.Int("port", cfg.DBPort),
	)
	defer func() {
		logger.Info("[AUTH-SERVICE] Cerrando conexiones de la base de datos...")
		_ = db.Close()
	}()

	if cfg.RunMigrations {
		logger.Info("[DATABASE] RUN_MIGRATIONS=true, ejecutando migraciones pendientes...")
		err := postgres.RunMigrations(db, "db/migrations")
		if err != nil {
			logger.Error("[AUTH-SERVICE] Error al ejecutar migraciones", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		logger.Debug("[DATABASE] RUN_MIGRATIONS=false, omitiendo migraciones")
	}

	// 4. Instanciar Adaptadores de Datos (Persistencia y Caché)
	logger.Debug("[AUTH-SERVICE] Instanciando repositorios de persistencia...")
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	logger.Debug("[AUTH-SERVICE] Conectando al cliente Redis...",
		slog.String("host", cfg.RedisHost),
		slog.Int("port", cfg.RedisPort),
		slog.Int("db", cfg.RedisDB),
	)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	redisCache := cache.NewRedisCache(redisClient)

	// 5. Instanciar Casos de Uso (Lógica de Aplicación)
	logger.Debug("[AUTH-SERVICE] Instanciando casos de uso...")
	registerUC := usecase.NewRegisterUseCase(userRepo, logger)
	loginUC := usecase.NewLoginUseCase(userRepo, redisCache, refreshTokenRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.RTExpiration, cfg.LoginMaxAttempts, cfg.LoginBlockMinutes, logger)
	refreshUC := usecase.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.RTExpiration, logger)
	logoutUC := usecase.NewLogoutUseCase(redisCache, refreshTokenRepo, logger)

	// 6. Configurar Middlewares
	jwtMiddleware := middleware.JWTAuth(cfg.JWTSecret, redisCache)
	loggingMiddleware := middleware.Logging

	// 7. Instanciar Router HTTP Único
	httpRouter := router.NewRouter(
		registerUC,
		loginUC,
		refreshUC,
		logoutUC,
		jwtMiddleware,
		loggingMiddleware,
		cfg.TimeoutRegister,
		cfg.TimeoutLogin,
		cfg.TimeoutMe,
		cfg.CookieSecure,
		cfg.RTExpiration,
	)

	// Calcular dinámicamente los límites globales del servidor HTTP basados en el timeout máximo de los endpoints
	maxTimeout := cfg.TimeoutRegister
	if cfg.TimeoutLogin > maxTimeout {
		maxTimeout = cfg.TimeoutLogin
	}
	if cfg.TimeoutMe > maxTimeout {
		maxTimeout = cfg.TimeoutMe
	}

	// Agregar un margen de seguridad de 2 segundos para dar tiempo al middleware de timeout a escribir la respuesta JSON controlada
	serverWriteTimeout := maxTimeout + 2*time.Second
	serverReadTimeout := maxTimeout + 2*time.Second

	// 8. Arrancar Servidor HTTP con Soporte de Apagado Elegante (Graceful Shutdown)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpRouter,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para escuchar errores de arranque
	serverErrors := make(chan error, 1)
	go func() {
		// Imprimir Banner de arranque en consola de forma limpia
		fmt.Printf(`
========================================================================
   ____    _   _  _____  _   _      ____   _____  ____  
  / ___|  | | | ||_   _|| | | |    / ___| |  ___||  _ \ 
 | |      | | | |  | |  | |_| |    \___ \ | |_   | |_) |
 | |___   | |_| |  | |  |  _  |     ___) ||  _|  |  _ < 
  \____|   \___/   |_|  |_| |_|    |____/ |_|    |_| \_\
                                                        
       HOMELAB COLLECTOR PLATFORM SKAOTICO V2 - AUTH SERVICE V1.3
========================================================================
 [+] PUERTO:          %s
 [+] SWAGGER UI:      http://localhost:%s/swagger/
 [+] SWAGGER JSON:    http://localhost:%s/swagger/doc.json

 [+] CONFIGURACIÓN DE SEGURIDAD:
      - JWT TTL        : %s
      - Refresh TTL    : %s
      - Max Login Fails: %d intentos
      - Block Window   : %d minutos

 [ENDPOINTS DISPONIBILIZADOS]:
  -> GET  /api/v1/health           [Estado del Servicio (Público)]
  -> POST /api/v1/auth/register    [Registro de usuarios (Público)] [Timeout: %s]
  -> POST /api/v1/auth/login       [Inicio de sesión / Token (Público)] [Timeout: %s]
  -> POST /api/v1/auth/refresh     [Rotación de Token (Público)] [Timeout: %s]
  -> GET  /api/v1/auth/me          [Visualizar claims del JWT (Protegido)] [Timeout: %s]
  -> POST /api/v1/auth/logout      [Invalidación y Blacklist (Protegido)] [Timeout: %s]
========================================================================
`, cfg.Port, cfg.Port, cfg.Port,
			cfg.JWTExpiration.String(), cfg.RTExpiration.String(), cfg.LoginMaxAttempts, cfg.LoginBlockMinutes,
			cfg.TimeoutRegister, cfg.TimeoutLogin, cfg.TimeoutLogin, cfg.TimeoutMe, cfg.TimeoutMe)

		logger.Info("[AUTH-SERVICE] Servidor HTTP escuchando",
			slog.String("addr", ":"+cfg.Port),
			slog.String("read_timeout", serverReadTimeout.String()),
			slog.String("write_timeout", serverWriteTimeout.String()),
		)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Esperar bloqueado por un error de arranque o una señal de apagado
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("[AUTH-SERVICE] Error crítico en el servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("[AUTH-SERVICE] Señal recibida. Iniciando apagado controlado...", slog.Any("signal", sig))

		// Otorgar un límite de 15 segundos para completar solicitudes HTTP activas
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("[AUTH-SERVICE] Error al apagar el servidor con gracia. Forzando cierre...", slog.Any("error", err))
			_ = server.Close()
		}
		logger.Info("[AUTH-SERVICE] Servidor HTTP apagado correctamente.")
	}
}
