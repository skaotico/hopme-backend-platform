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

	"c4-sensor/internal/infra/config"
	"c4-sensor/internal/infra/database"
	"c4-sensor/internal/infra/database/postgres"
	"c4-sensor/internal/infra/http/middleware"
	"c4-sensor/internal/infra/http/router"
	"c4-sensor/internal/infra/observability"
	"c4-sensor/internal/usecase"
)

// Run centraliza toda la inicialización, inyección de dependencias y arranque del servicio
func Run() {
	cfg := config.Load()

	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "sensor-service",
		ServiceVersion: "v1.0",
		Environment:    cfg.Environment,
	})

	logger.Info("[SENSOR-SERVICE] Iniciando servicio...",
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
		logger.Error("[SENSOR-SERVICE] Error crítico al conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("[SENSOR-SERVICE] Cerrando conexiones de la base de datos...")
		if err := db.Close(); err != nil {
			logger.Error("[SENSOR-SERVICE] Error al cerrar conexión de base de datos", slog.Any("error", err))
		}
	}()

	if cfg.RunMigrations {
		logger.Info("[DATABASE] RUN_MIGRATIONS=true, ejecutando migraciones pendientes...")
		err := postgres.RunMigrations(db, "db/migrations")
		if err != nil {
			logger.Error("[SENSOR-SERVICE] Error al ejecutar migraciones", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// Repositorios
	sensorRepo := postgres.NewSensorRepository(db, logger)
	lecturaRepo := postgres.NewLecturaRepository(db, logger)
	alertaRepo := postgres.NewAlertaRepository(db, logger)

	// Casos de Uso
	sensorUC := usecase.NewSensorUseCase(sensorRepo, logger)
	lecturaUC := usecase.NewLecturaUseCase(lecturaRepo, sensorRepo, logger)
	alertaUC := usecase.NewAlertaUseCase(alertaRepo, sensorRepo, logger)

	// Middlewares
	loggingMid := middleware.Logging

	// Router
	r := router.NewRouter(sensorUC, lecturaUC, alertaUC, logger, loggingMid)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para recibir señales del sistema operativo para apagado ordenado
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf(`
========================================================================
  ____  _____ _   _ ____   ___  ____  
 / ___|| ____| \ | / ___| / _ \|  _ \ 
 \___ \|  _| |  \| \___ \| | | | |_) |
  ___) | |___| |\  |___) | |_| |  _ < 
 |____/|_____|_| \_|____/ \___/|_| \_\
                                                    
         HOMELAB COLLECTOR PLATFORM - SENSOR SERVICE V1.0
========================================================================
 [+] PUERTO:          %s
========================================================================
`, cfg.Port)
		logger.Info("[SENSOR-SERVICE] Servidor HTTP escuchando", slog.String("addr", ":"+cfg.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("[SENSOR-SERVICE] Error al arrancar el servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-stop

	logger.Info("[SENSOR-SERVICE] Recibida señal de parada, iniciando apagado ordenado (graceful shutdown)...")

	// Contexto con timeout de 10 segundos para dar tiempo a procesar peticiones en vuelo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("[SENSOR-SERVICE] Error durante el apagado del servidor", slog.Any("error", err))
	} else {
		logger.Info("[SENSOR-SERVICE] Servidor HTTP detenido correctamente.")
	}
}
