package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"c4-sensor/internal/domain/model"
	"c4-sensor/internal/domain/port"
	"github.com/google/uuid"
)

type lecturaUseCase struct {
	repo       port.LecturaRepository
	sensorRepo port.SensorRepository
	log        *slog.Logger
}

// NewLecturaUseCase crea una nueva instancia del caso de uso de lecturas.
func NewLecturaUseCase(repo port.LecturaRepository, sensorRepo port.SensorRepository, log *slog.Logger) port.LecturaUseCase {
	return &lecturaUseCase{
		repo:       repo,
		sensorRepo: sensorRepo,
		log:        log.With(slog.String("component", "lectura_usecase")),
	}
}

func (uc *lecturaUseCase) Create(ctx context.Context, sensorID uuid.UUID, req model.CreateLecturaRequest) (*model.LecturaSensor, error) {
	log := uc.log.With(slog.String("operation", "Create"), slog.String("sensor_id", sensorID.String()))
	log.InfoContext(ctx, "registrando nueva lectura de sensor")

	_, err := uc.sensorRepo.GetByID(ctx, sensorID)
	if err != nil {
		log.WarnContext(ctx, "sensor no existe", slog.Any("error", err))
		return nil, errors.New("sensor no encontrado")
	}

	var lecturaTime time.Time
	if req.FechaLectura != nil {
		lecturaTime = *req.FechaLectura
	} else {
		lecturaTime = time.Now().UTC()
	}

	lectura := &model.LecturaSensor{
		ID:                uuid.New(),
		SensorID:          sensorID,
		Valor:             req.Valor,
		FechaLectura:      lecturaTime,
		BateriaPorcentaje: req.BateriaPorcentaje,
		Observacion:       req.Observacion,
	}

	if err := uc.repo.Create(ctx, lectura); err != nil {
		log.ErrorContext(ctx, "error al registrar lectura en repositorio", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "lectura de sensor registrada exitosamente", slog.String("lectura_id", lectura.ID.String()))
	return lectura, nil
}

func (uc *lecturaUseCase) ListBySensorID(ctx context.Context, sensorID uuid.UUID) ([]model.LecturaSensor, error) {
	log := uc.log.With(slog.String("operation", "ListBySensorID"), slog.String("sensor_id", sensorID.String()))
	log.DebugContext(ctx, "consultando lecturas por sensor_id")

	_, err := uc.sensorRepo.GetByID(ctx, sensorID)
	if err != nil {
		log.WarnContext(ctx, "sensor no existe", slog.Any("error", err))
		return nil, errors.New("sensor no encontrado")
	}

	lecturas, err := uc.repo.ListBySensorID(ctx, sensorID)
	if err != nil {
		log.ErrorContext(ctx, "error al obtener lecturas", slog.Any("error", err))
		return nil, err
	}

	return lecturas, nil
}
