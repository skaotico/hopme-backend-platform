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

	"c4-arbol/internal/infra/config"
	"c4-arbol/internal/infra/database"
	"c4-arbol/internal/infra/database/postgres"
	"c4-arbol/internal/infra/http/middleware"
	"c4-arbol/internal/infra/http/router"
	"c4-arbol/internal/infra/observability"
	"c4-arbol/internal/usecase"
)

// Run centraliza toda la inicialización, inyección de dependencias y arranque del servicio
func Run() {
	cfg := config.Load()

	logger := observability.New(observability.Config{
		Level:          observability.LogLevel(cfg.LogLevel),
		ServiceName:    "arbol-service",
		ServiceVersion: "v1.0",
		Environment:    cfg.Environment,
	})

	logger.Info("[ARBOL-SERVICE] Iniciando servicio...",
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
		logger.Error("[ARBOL-SERVICE] Error crítico al conectar a la base de datos", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("[ARBOL-SERVICE] Error al cerrar conexión de base de datos", slog.Any("error", err))
		}
	}()

	// Repositorios
	arbolRepo := postgres.NewArbolRepository(db, logger)
	historialRepo := postgres.NewHistorialRepository(db, logger)
	medicionRepo := postgres.NewMedicionRepository(db, logger)

	// Casos de Uso
	arbolUC := usecase.NewArbolUseCase(arbolRepo, logger)
	historialUC := usecase.NewHistorialUseCase(historialRepo, arbolRepo, logger)
	medicionUC := usecase.NewMedicionUseCase(medicionRepo, arbolRepo, logger)

	// Middlewares
	loggingMid := middleware.Logging

	// Router
	r := router.NewRouter(arbolUC, historialUC, medicionUC, loggingMid)

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
   ___   ____   ____    ___   _      
  / _ \ |  _ \ | __ )  / _ \ | |     
 / /_\ \| |_) ||  _ \ | | | || |     
/ /___ ||  _ < | |_) || |_| || |___  
/_/   \_\|_| \_\|____/ \___/ |_____| 
                                                       
       HOMELAB COLLECTOR PLATFORM SKAOTICO V2 - ARBOL SERVICE V1.0
========================================================================
 [+] PUERTO:          %s
 [+] SWAGGER UI:      http://localhost:%s/swagger/
 [+] SWAGGER JSON:    http://localhost:%s/swagger/doc.json

 [ENDPOINTS DISPONIBILIZADOS]:
  -> GET    /api/v1/health             [Estado del Servicio (Público)]
  -> POST   /api/v1/arboles            [Registrar árbol]
  -> GET    /api/v1/arboles            [Listar árboles]
  -> GET    /api/v1/arboles/{id}       [Obtener árbol por ID]
  -> PUT    /api/v1/arboles/{id}       [Actualizar árbol]
  -> DELETE /api/v1/arboles/{id}       [Eliminar árbol]
  -> POST   /api/v1/arboles/{id}/historial [Registrar cambio de estado]
  -> GET    /api/v1/arboles/{id}/historial [Listar historial de estados]
  -> POST   /api/v1/arboles/{id}/mediciones [Registrar medición dasométrica]
  -> GET    /api/v1/arboles/{id}/mediciones [Listar mediciones de un árbol]
========================================================================
`, cfg.Port, cfg.Port, cfg.Port)

		logger.Info("[ARBOL-SERVICE] Servidor HTTP escuchando",
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
			logger.Error("[ARBOL-SERVICE] Error crítico en el servidor HTTP", slog.Any("error", err))
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("[ARBOL-SERVICE] Señal recibida. Iniciando apagado controlado...", slog.Any("signal", sig))

		// Otorgar un límite de 15 segundos para completar solicitudes HTTP activas
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("[ARBOL-SERVICE] Error al apagar el servidor con gracia. Forzando cierre...", slog.Any("error", err))
			_ = server.Close()
		}
		logger.Info("[ARBOL-SERVICE] Servidor HTTP apagado correctamente.")
	}
}
