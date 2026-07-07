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

type alertaUseCase struct {
	repo       port.AlertaRepository
	sensorRepo port.SensorRepository
	log        *slog.Logger
}

// NewAlertaUseCase crea una nueva instancia del caso de uso de alertas.
func NewAlertaUseCase(repo port.AlertaRepository, sensorRepo port.SensorRepository, log *slog.Logger) port.AlertaUseCase {
	return &alertaUseCase{
		repo:       repo,
		sensorRepo: sensorRepo,
		log:        log.With(slog.String("component", "alerta_usecase")),
	}
}

func (uc *alertaUseCase) Create(ctx context.Context, sensorID uuid.UUID, req model.CreateAlertaRequest) (*model.Alerta, error) {
	log := uc.log.With(slog.String("operation", "Create"), slog.String("sensor_id", sensorID.String()))
	log.InfoContext(ctx, "creando nueva alerta")

	_, err := uc.sensorRepo.GetByID(ctx, sensorID)
	if err != nil {
		log.WarnContext(ctx, "sensor no existe", slog.Any("error", err))
		return nil, errors.New("sensor no encontrado")
	}

	var alertaTime time.Time
	if req.FechaAlerta != nil {
		alertaTime = *req.FechaAlerta
	} else {
		alertaTime = time.Now().UTC()
	}

	alerta := &model.Alerta{
		ID:             uuid.New(),
		SensorID:       sensorID,
		Tipo:           req.Tipo,
		ValorDetectado: req.ValorDetectado,
		Umbral:         req.Umbral,
		FechaAlerta:    alertaTime,
		Estado:         "activa",
		Observacion:    req.Observacion,
	}

	if err := uc.repo.Create(ctx, alerta); err != nil {
		log.ErrorContext(ctx, "error al registrar alerta en repositorio", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "alerta registrada exitosamente", slog.String("alerta_id", alerta.ID.String()))
	return alerta, nil
}

func (uc *alertaUseCase) List(ctx context.Context, sensorID *uuid.UUID, estado *string) ([]model.Alerta, error) {
	log := uc.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "consultando lista de alertas")

	if sensorID != nil {
		_, err := uc.sensorRepo.GetByID(ctx, *sensorID)
		if err != nil {
			log.WarnContext(ctx, "sensor no existe", slog.Any("error", err))
			return nil, errors.New("sensor no encontrado")
		}
	}

	alertas, err := uc.repo.List(ctx, sensorID, estado)
	if err != nil {
		log.ErrorContext(ctx, "error al consultar alertas en repositorio", slog.Any("error", err))
		return nil, err
	}

	return alertas, nil
}

func (uc *alertaUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, req model.UpdateAlertaStatusRequest) (*model.Alerta, error) {
	log := uc.log.With(slog.String("operation", "UpdateStatus"), slog.String("alerta_id", id.String()))
	log.InfoContext(ctx, "actualizando estado de alerta")

	alerta, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, "alerta no encontrada", slog.Any("error", err))
		return nil, err
	}

	alerta.Estado = req.Estado
	if req.Observacion != nil {
		alerta.Observacion = req.Observacion
	}

	if err := uc.repo.Update(ctx, alerta); err != nil {
		log.ErrorContext(ctx, "error al actualizar alerta en repositorio", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "estado de alerta actualizado exitosamente")
	return alerta, nil
}
