package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-zona/internal/domain/model"

	"github.com/google/uuid"
)

// Update actualiza los campos modificables de una Zona existente.
func (uc *zonaUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateZonaRequest) (*model.Zona, error) {
	log := uc.log.With(slog.String("operation", "Update"), slog.String("zona_id", id.String()))
	log.InfoContext(ctx, "iniciando actualización de zona")

	zona, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.WarnContext(ctx, "zona no encontrada para actualizar",
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.EcoparqueID != nil {
		zona.EcoparqueID = *req.EcoparqueID
	}
	if req.Nombre != nil {
		zona.Nombre = *req.Nombre
	}
	if req.Descripcion != nil {
		zona.Descripcion = req.Descripcion
	}
	if req.AreaM2 != nil {
		zona.AreaM2 = req.AreaM2
	}

	now := time.Now().UTC()
	zona.FechaActualizacion = &now

	if err = uc.repo.Update(ctx, zona); err != nil {
		log.ErrorContext(ctx, "error al actualizar zona en repositorio",
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "zona actualizada exitosamente")
	return zona, nil
}
