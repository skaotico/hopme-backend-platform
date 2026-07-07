package usecase

import (
	"context"
	"log/slog"

	"c4-zona/internal/domain/model"

	"github.com/google/uuid"
)

// GetByID obtiene una Zona por su identificador único.
func (uc *zonaUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error) {
	log := uc.log.With(slog.String("operation", "GetByID"), slog.String("zona_id", id.String()))
	log.InfoContext(ctx, "buscando zona por ID")

	zona, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.WarnContext(ctx, "zona no encontrada",
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "zona encontrada")
	return zona, nil
}
