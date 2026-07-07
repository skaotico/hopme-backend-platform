package usecase

import (
	"log/slog"

	"c4-sensor/internal/domain/port"
)

type sensorUseCase struct {
	repo port.SensorRepository
	log  *slog.Logger
}

// NewSensorUseCase crea una nueva instancia del caso de uso de Sensores.
func NewSensorUseCase(repo port.SensorRepository, log *slog.Logger) port.SensorUseCase {
	return &sensorUseCase{
		repo: repo,
		log:  log.With(slog.String("component", "sensor_usecase")),
	}
}
