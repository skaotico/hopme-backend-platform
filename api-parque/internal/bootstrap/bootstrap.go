package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"c4-parque/internal/infra/config"
	"c4-parque/internal/infra/database"
	"c4-parque/internal/infra/database/postgres"
	"c4-parque/internal/infra/http/middleware"
	"c4-parque/internal/infra/http/router"
	"c4-parque/internal/infra/observability"
	"c4-parque/internal/usecase"
)

// Run centraliza toda la inicialización, inyección de dependencias y arranque del servicio
func Run() {
	cfg := config.Load()

	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "parque-service",
		ServiceVersion: "v1.0",
		Environment:    cfg.Environment,
	})

	logger.Info("[PARQUE-SERVICE] Iniciando servicio...",
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
		logger.Error("[PARQUE-SERVICE] Error crítico al conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()



	// Dependencias
	ecoparqueRepo := postgres.NewEcoparqueRepository(db)
	ecoparqueUC := usecase.NewEcoparqueUseCase(ecoparqueRepo)

	// Middlewares
	loggingMid := middleware.Logging

	// Router
	r := router.NewRouter(ecoparqueUC, loggingMid)
	

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("\n")
		fmt.Printf("===========================================================\n")
		fmt.Printf("  🚀 PARQUE-SERVICE IS RUNNING! [%s]\n", cfg.Environment)
		fmt.Printf("===========================================================\n")
		fmt.Printf("  => API Base:   http://localhost:%s/api/v1\n", cfg.Port)
		fmt.Printf("  => Health:     http://localhost:%s/api/v1/health\n", cfg.Port)
		fmt.Printf("  => Swagger UI: http://localhost:%s/swagger/index.html\n", cfg.Port)
		fmt.Printf("===========================================================\n\n")

		logger.Info("[PARQUE-SERVICE] Servidor escuchando", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("[PARQUE-SERVICE] Error en servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("[PARQUE-SERVICE] Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("[PARQUE-SERVICE] Apagado forzado del servidor", slog.Any("error", err))
	}
	logger.Info("[PARQUE-SERVICE] Servidor detenido correctamente")
}
