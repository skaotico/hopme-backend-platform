package usecase

import (
	"context"
	"log/slog"

	"c4-sensor/internal/domain/model"
)

func (uc *sensorUseCase) List(ctx context.Context, filter model.SensorFilter) ([]model.Sensor, error) {
	log := uc.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "iniciando consulta de lista de sensores con filtros")

	sensores, err := uc.repo.List(ctx, filter)
	if err != nil {
		log.ErrorContext(ctx, "error al consultar lista de sensores", slog.Any("error", err))
		return nil, err
	}

	log.DebugContext(ctx, "consulta de sensores completada", slog.Int("count", len(sensores)))
	return sensores, nil
}
