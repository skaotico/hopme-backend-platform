package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (uc *sensorUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	log := uc.log.With(slog.String("operation", "Delete"), slog.String("sensor_id", id.String()))
	log.InfoContext(ctx, "iniciando eliminación de sensor")

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "sensor no encontrado para eliminar", slog.Any("error", err))
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		log.ErrorContext(ctx, "error al eliminar sensor en repositorio", slog.Any("error", err))
		return err
	}

	log.InfoContext(ctx, "sensor eliminado exitosamente")
	return nil
}
