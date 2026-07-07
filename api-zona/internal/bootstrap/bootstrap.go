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

	// Canal para escuchar errores de arranque
	serverErrors := make(chan error, 1)
	go func() {
		// Imprimir Banner de arranque en consola de forma limpia
		fmt.Printf(`
========================================================================
  _______  ___   _   _  ___      ____   _____  ____  
 |__  /  \/   \ | \ | |/ _ \    / ___| |  ___||  _ \ 
   / /| |  | | ||  \| | | | |   \___ \ | |_   | |_) |
  / /_| |  |_| || |\  | |_| |    ___) ||  _|  |  _ < 
 /____|___\___/ |_| \_|\___/    |____/ |_|    |_| \_\
                                                      
       HOMELAB COLLECTOR PLATFORM SKAOTICO V2 - ZONA SERVICE V1.0
========================================================================
 [+] PUERTO:          %s
 [+] SWAGGER UI:      http://localhost:%s/swagger/
 [+] SWAGGER JSON:    http://localhost:%s/swagger/doc.json

 [ENDPOINTS DISPONIBILIZADOS]:
  -> GET    /api/v1/health          [Estado del Servicio (Público)]
  -> POST   /api/v1/zonas           [Crear zona]
  -> GET    /api/v1/zonas           [Listar zonas]
  -> GET    /api/v1/zonas/{id}      [Obtener zona por ID]
  -> PUT    /api/v1/zonas/{id}      [Actualizar zona]
  -> DELETE /api/v1/zonas/{id}      [Eliminar zona]
========================================================================
`, cfg.Port, cfg.Port, cfg.Port)

		logger.Info("[ZONA-SERVICE] Servidor HTTP escuchando",
			slog.String("addr", ":"+cfg.Port),
			slog.String("read_timeout", server.ReadTimeout.String()),
			slog.String("write_timeout", server.WriteTimeout.String()),
		)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Esperar bloqueado por un error de arranque o una señal de apagado
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("[ZONA-SERVICE] Error crítico en el servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("[ZONA-SERVICE] Señal recibida. Iniciando apagado controlado...", slog.Any("signal", sig))

		// Otorgar un límite de 15 segundos para completar solicitudes HTTP activas
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("[ZONA-SERVICE] Error al apagar el servidor con gracia. Forzando cierre...", slog.Any("error", err))
			_ = server.Close()
		}
		logger.Info("[ZONA-SERVICE] Servidor HTTP apagado correctamente.")
	}
}
