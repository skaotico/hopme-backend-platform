package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-sensor/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *sensorUseCase) Create(ctx context.Context, req model.CreateSensorRequest) (*model.Sensor, error) {
	log := uc.log.With(slog.String("operation", "Create"))
	log.InfoContext(ctx, "iniciando registro de sensor")

	var installTime time.Time
	if req.FechaInstalacion != nil {
		installTime = *req.FechaInstalacion
	} else {
		installTime = time.Now().UTC()
	}

	activo := true
	sensor := &model.Sensor{
		ID:               uuid.New(),
		Codigo:           req.Codigo,
		TipoSensorID:     req.TipoSensorID,
		ArbolID:          req.ArbolID,
		EstanqueID:       req.EstanqueID,
		Fabricante:       req.Fabricante,
		Modelo:           req.Modelo,
		FechaInstalacion: &installTime,
		Activo:           &activo,
	}

	if err := uc.repo.Create(ctx, sensor); err != nil {
		log.ErrorContext(ctx, "error al crear sensor en repositorio",
			slog.String("sensor_id", sensor.ID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "sensor registrado exitosamente",
		slog.String("sensor_id", sensor.ID.String()),
	)
	return sensor, nil
}
