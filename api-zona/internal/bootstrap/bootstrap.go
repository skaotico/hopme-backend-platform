package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"c4-zona/internal/infra/config"
	"c4-zona/internal/infra/database"
	"c4-zona/internal/infra/database/postgres"
	"c4-zona/internal/infra/http/middleware"
	"c4-zona/internal/infra/http/router"
	"c4-zona/internal/infra/observability"
	"c4-zona/internal/usecase"
)

// Run centraliza toda la inicialización, inyección de dependencias y arranque del servicio
func Run() {
	cfg := config.Load()

	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "zona-service",
		ServiceVersion: "v1.0",
		Environment:    cfg.Environment,
	})

	logger.Info("[ZONA-SERVICE] Iniciando servicio...",
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
		logger.Error("[ZONA-SERVICE] Error crítico al conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("[ZONA-SERVICE] Error al cerrar conexión de base de datos", slog.Any("error", err))
		}
	}()



	// Dependencias
	zonaRepo := postgres.NewZonaRepository(db)
	zonaUC := usecase.NewZonaUseCase(zonaRepo)

	// Middlewares
	loggingMid := middleware.Logging

	// Router
	r := router.NewRouter(zonaUC, loggingMid)
	

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("[ZONA-SERVICE] Servidor escuchando", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("[ZONA-SERVICE] Error en servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("[ZONA-SERVICE] Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("[ZONA-SERVICE] Apagado forzado del servidor", slog.Any("error", err))
	}
	logger.Info("[ZONA-SERVICE] Servidor detenido correctamente")
}
