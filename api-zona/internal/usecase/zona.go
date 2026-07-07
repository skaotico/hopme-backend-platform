package usecase

import (
	"log/slog"

	"c4-zona/internal/domain/port"
)

type zonaUseCase struct {
	repo port.ZonaRepository
	log  *slog.Logger
}

// NewZonaUseCase crea una nueva instancia del caso de uso de Zona.
func NewZonaUseCase(repo port.ZonaRepository, log *slog.Logger) port.ZonaUseCase {
	return &zonaUseCase{
		repo: repo,
		log:  log.With(slog.String("component", "zona_usecase")),
	}
}
