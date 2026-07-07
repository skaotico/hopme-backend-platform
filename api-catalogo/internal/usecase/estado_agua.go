package usecase

import (
	"context"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type estadoAguaUseCase struct {
	repo   port.EstadoAguaRepository
	logger *slog.Logger
}

// NewEstadoAguaUseCase crea una instancia del caso de uso de estados de agua
func NewEstadoAguaUseCase(repo port.EstadoAguaRepository, logger *slog.Logger) port.EstadoAguaUseCase {
	return &estadoAguaUseCase{
		repo:   repo,
		logger: observability.WithComponent(logger, "usecase.estado_agua"),
	}
}

func (uc *estadoAguaUseCase) Create(ctx context.Context, req model.CreateEstadoAguaRequest) (*model.EstadoAgua, error) {
	uc.logger.Debug("[ESTADO_AGUA] Creando nuevo estado", slog.String("codigo", req.Codigo))
	e := &model.EstadoAgua{ID: uuid.New(), Codigo: req.Codigo, Nombre: req.Nombre}
	if err := uc.repo.Create(ctx, e); err != nil {
		uc.logger.Error("[ESTADO_AGUA] Error al crear estado de agua", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_AGUA] Estado de agua creado", slog.String("id", e.ID.String()))
	return e, nil
}

func (uc *estadoAguaUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoAgua, error) {
	uc.logger.Debug("[ESTADO_AGUA] Buscando estado de agua por ID", slog.String("id", id.String()))
	res, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("[ESTADO_AGUA] Error al buscar estado de agua por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if res == nil {
		uc.logger.Warn("[ESTADO_AGUA] Estado de agua no encontrado por ID", slog.String("id", id.String()))
		return nil, nil
	}
	uc.logger.Info("[ESTADO_AGUA] Estado de agua encontrado por ID", slog.String("id", id.String()))
	return res, nil
}

func (uc *estadoAguaUseCase) List(ctx context.Context) ([]model.EstadoAgua, error) {
	uc.logger.Debug("[ESTADO_AGUA] Listando estados de agua")
	res, err := uc.repo.List(ctx)
	if err != nil {
		uc.logger.Error("[ESTADO_AGUA] Error al listar estados de agua", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_AGUA] Estados de agua listados con éxito", slog.Int("cantidad", len(res)))
	return res, nil
}

func (uc *estadoAguaUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoAguaRequest) (*model.EstadoAgua, error) {
	uc.logger.Debug("[ESTADO_AGUA] Actualizando estado de agua", slog.String("id", id.String()))
	updated, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("[ESTADO_AGUA] Error al actualizar estado de agua", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_AGUA] Estado de agua actualizado", slog.String("id", id.String()))
	return updated, nil
}

func (uc *estadoAguaUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Debug("[ESTADO_AGUA] Eliminando estado de agua", slog.String("id", id.String()))
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("[ESTADO_AGUA] Error al eliminar estado de agua", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	uc.logger.Info("[ESTADO_AGUA] Estado de agua eliminado con éxito", slog.String("id", id.String()))
	return nil
}
