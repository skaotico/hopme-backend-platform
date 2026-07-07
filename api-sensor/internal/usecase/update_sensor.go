package usecase

import (
	"context"
	"log/slog"

	"c4-sensor/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *sensorUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateSensorRequest) (*model.Sensor, error) {
	log := uc.log.With(slog.String("operation", "Update"), slog.String("sensor_id", id.String()))
	log.InfoContext(ctx, "iniciando actualización de sensor")

	sensor, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "sensor no encontrado para actualizar", slog.Any("error", err))
		return nil, err
	}

	if req.Codigo != nil {
		sensor.Codigo = req.Codigo
	}
	if req.TipoSensorID != nil {
		sensor.TipoSensorID = *req.TipoSensorID
	}
	if req.ArbolID != nil {
		sensor.ArbolID = req.ArbolID
	}
	if req.EstanqueID != nil {
		sensor.EstanqueID = req.EstanqueID
	}
	if req.Fabricante != nil {
		sensor.Fabricante = req.Fabricante
	}
	if req.Modelo != nil {
		sensor.Modelo = req.Modelo
	}
	if req.FechaInstalacion != nil {
		sensor.FechaInstalacion = req.FechaInstalacion
	}
	if req.Activo != nil {
		sensor.Activo = req.Activo
	}

	if err := uc.repo.Update(ctx, sensor); err != nil {
		log.ErrorContext(ctx, "error al actualizar sensor en repositorio", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "sensor actualizado exitosamente")
	return sensor, nil
}
