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

	"c4-catalogo/internal/infra/config"
	"c4-catalogo/internal/infra/database"
	"c4-catalogo/internal/infra/database/postgres"
	catalogoHTTP "c4-catalogo/internal/infra/http"
	"c4-catalogo/internal/infra/http/middleware"
	"c4-catalogo/internal/infra/http/router"
	"c4-catalogo/internal/infra/observability"
	"c4-catalogo/internal/usecase"
)

// Run centraliza la inicialización y el arranque del servicio Catalogo
func Run() {
	cfg := config.Load()

	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "catalogo-service",
		ServiceVersion: "v1.0",
		Environment:    cfg.Environment,
	})

	logger.Info("[CATALOGO-SERVICE] Iniciando servicio de catálogo...",
		slog.String("log_level", cfg.LogLevel),
		slog.String("environment", cfg.Environment),
		slog.String("port", cfg.Port),
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
		logger.Error("[CATALOGO-SERVICE] Error crítico al conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("[DATABASE] Conexión a PostgreSQL establecida correctamente")
	defer func() {
		logger.Info("[CATALOGO-SERVICE] Cerrando conexiones de la base de datos...")
		_ = db.Close()
	}()

	if cfg.RunMigrations {
		logger.Info("[DATABASE] RUN_MIGRATIONS=true, ejecutando migraciones pendientes...")
		err := postgres.RunMigrations(db, "db/migrations")
		if err != nil {
			logger.Error("[CATALOGO-SERVICE] Error al ejecutar migraciones", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// Instanciar repositorios
	especieRepo := postgres.NewEspecieArbolRepository(db, logger)
	estadoArbolRepo := postgres.NewEstadoArbolRepository(db, logger)
	estadoEstanqueRepo := postgres.NewEstadoEstanqueRepository(db, logger)
	estadoAguaRepo := postgres.NewEstadoAguaRepository(db, logger)
	tipoSensorRepo := postgres.NewTipoSensorRepository(db, logger)

	// Instanciar casos de uso
	especieUC := usecase.NewEspecieArbolUseCase(especieRepo, logger)
	estadoArbolUC := usecase.NewEstadoArbolUseCase(estadoArbolRepo, logger)
	estadoEstanqueUC := usecase.NewEstadoEstanqueUseCase(estadoEstanqueRepo, logger)
	estadoAguaUC := usecase.NewEstadoAguaUseCase(estadoAguaRepo, logger)
	tipoSensorUC := usecase.NewTipoSensorUseCase(tipoSensorRepo, logger)

	// Instanciar handlers
	especieH := catalogoHTTP.NewEspecieArbolHandler(especieUC)
	estadoArbolH := catalogoHTTP.NewEstadoArbolHandler(estadoArbolUC)
	estadoEstanqueH := catalogoHTTP.NewEstadoEstanqueHandler(estadoEstanqueUC)
	estadoAguaH := catalogoHTTP.NewEstadoAguaHandler(estadoAguaUC)
	tipoSensorH := catalogoHTTP.NewTipoSensorHandler(tipoSensorUC)

	// Router
	httpRouter := router.NewRouter(
		especieH,
		estadoArbolH,
		estadoEstanqueH,
		estadoAguaH,
		tipoSensorH,
		middleware.Logging,
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		fmt.Printf(`
========================================================================
   ____    _   _____   _   _      ___   ____   ___  
  / ___|  / \ |_   _| / \ | |    / _ \ / ___| / _ \ 
 | |     / _ \  | |  / _ \| |   | | | | |  _ | | | |
 | |___ / ___ \ | | / ___ \ |___| |_| | |_| || |_| |
  \____/_/   \_\|_|/_/   \_\_____|\___/ \____| \___/
                                                    
        HOMELAB COLLECTOR PLATFORM - CATALOGO SERVICE V1.0
========================================================================
 [+] PUERTO:          %s
========================================================================
`, cfg.Port)
		logger.Info("[CATALOGO-SERVICE] Servidor HTTP escuchando", slog.String("addr", ":"+cfg.Port))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("[CATALOGO-SERVICE] Error crítico en el servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("[CATALOGO-SERVICE] Señal recibida. Iniciando apagado controlado...", slog.Any("signal", sig))
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("[CATALOGO-SERVICE] Error al apagar el servidor con gracia. Forzando cierre...", slog.Any("error", err))
			_ = server.Close()
		}
		logger.Info("[CATALOGO-SERVICE] Servidor HTTP apagado correctamente.")
	}
}
