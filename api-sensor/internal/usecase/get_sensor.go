package usecase

import (
	"context"
	"log/slog"

	"c4-sensor/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *sensorUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Sensor, error) {
	log := uc.log.With(slog.String("operation", "GetByID"), slog.String("sensor_id", id.String()))
	log.DebugContext(ctx, "buscando sensor por ID")

	sensor, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "sensor no encontrado en base de datos", slog.Any("error", err))
		return nil, err
	}

	return sensor, nil
}
