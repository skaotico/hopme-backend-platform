package usecase

import (
	"context"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type tipoSensorUseCase struct {
	repo   port.TipoSensorRepository
	logger *slog.Logger
}

// NewTipoSensorUseCase crea una instancia del caso de uso de tipos de sensor
func NewTipoSensorUseCase(repo port.TipoSensorRepository, logger *slog.Logger) port.TipoSensorUseCase {
	return &tipoSensorUseCase{
		repo:   repo,
		logger: observability.WithComponent(logger, "usecase.tipo_sensor"),
	}
}

func (uc *tipoSensorUseCase) Create(ctx context.Context, req model.CreateTipoSensorRequest) (*model.TipoSensor, error) {
	uc.logger.Debug("[TIPO_SENSOR] Creando nuevo tipo de sensor", slog.String("codigo", req.Codigo))
	e := &model.TipoSensor{
		ID:           uuid.New(),
		Codigo:       req.Codigo,
		Nombre:       req.Nombre,
		UnidadMedida: req.UnidadMedida,
		Descripcion:  req.Descripcion,
	}
	if err := uc.repo.Create(ctx, e); err != nil {
		uc.logger.Error("[TIPO_SENSOR] Error al crear tipo de sensor", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[TIPO_SENSOR] Tipo de sensor creado", slog.String("id", e.ID.String()))
	return e, nil
}

func (uc *tipoSensorUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.TipoSensor, error) {
	uc.logger.Debug("[TIPO_SENSOR] Buscando tipo de sensor por ID", slog.String("id", id.String()))
	res, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("[TIPO_SENSOR] Error al buscar tipo de sensor por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if res == nil {
		uc.logger.Warn("[TIPO_SENSOR] Tipo de sensor no encontrado por ID", slog.String("id", id.String()))
		return nil, nil
	}
	uc.logger.Info("[TIPO_SENSOR] Tipo de sensor encontrado por ID", slog.String("id", id.String()))
	return res, nil
}

func (uc *tipoSensorUseCase) List(ctx context.Context) ([]model.TipoSensor, error) {
	uc.logger.Debug("[TIPO_SENSOR] Listando tipos de sensor")
	res, err := uc.repo.List(ctx)
	if err != nil {
		uc.logger.Error("[TIPO_SENSOR] Error al listar tipos de sensor", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[TIPO_SENSOR] Tipos de sensor listados con éxito", slog.Int("cantidad", len(res)))
	return res, nil
}

func (uc *tipoSensorUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateTipoSensorRequest) (*model.TipoSensor, error) {
	uc.logger.Debug("[TIPO_SENSOR] Actualizando tipo de sensor", slog.String("id", id.String()))
	updated, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("[TIPO_SENSOR] Error al actualizar tipo de sensor", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[TIPO_SENSOR] Tipo de sensor actualizado", slog.String("id", id.String()))
	return updated, nil
}

func (uc *tipoSensorUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Debug("[TIPO_SENSOR] Eliminando tipo de sensor", slog.String("id", id.String()))
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("[TIPO_SENSOR] Error al eliminar tipo de sensor", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	uc.logger.Info("[TIPO_SENSOR] Tipo de sensor eliminado con éxito", slog.String("id", id.String()))
	return nil
}
