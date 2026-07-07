package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-zona/internal/domain/model"

	"github.com/google/uuid"
)

// Create crea una nueva Zona y la persiste en el repositorio.
func (uc *zonaUseCase) Create(ctx context.Context, req model.CreateZonaRequest) (*model.Zona, error) {
	log := uc.log.With(slog.String("operation", "Create"))
	log.InfoContext(ctx, "iniciando creación de zona",
		slog.String("ecoparque_id", req.EcoparqueID.String()),
		slog.String("nombre", req.Nombre),
	)

	zona := &model.Zona{
		ID:            uuid.New(),
		EcoparqueID:   req.EcoparqueID,
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		AreaM2:        req.AreaM2,
		FechaCreacion: time.Now().UTC(),
	}

	if err := uc.repo.Create(ctx, zona); err != nil {
		log.ErrorContext(ctx, "error al crear zona en repositorio",
			slog.String("zona_id", zona.ID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "zona creada exitosamente",
		slog.String("zona_id", zona.ID.String()),
	)
	return zona, nil
}
