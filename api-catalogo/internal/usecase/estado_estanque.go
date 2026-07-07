package usecase

import (
	"context"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type estadoEstanqueUseCase struct {
	repo   port.EstadoEstanqueRepository
	logger *slog.Logger
}

// NewEstadoEstanqueUseCase crea una instancia del caso de uso de estados de estanque
func NewEstadoEstanqueUseCase(repo port.EstadoEstanqueRepository, logger *slog.Logger) port.EstadoEstanqueUseCase {
	return &estadoEstanqueUseCase{
		repo:   repo,
		logger: observability.WithComponent(logger, "usecase.estado_estanque"),
	}
}

func (uc *estadoEstanqueUseCase) Create(ctx context.Context, req model.CreateEstadoEstanqueRequest) (*model.EstadoEstanque, error) {
	uc.logger.Debug("[ESTADO_ESTANQUE] Creando nuevo estado", slog.String("codigo", req.Codigo))
	e := &model.EstadoEstanque{ID: uuid.New(), Codigo: req.Codigo, Nombre: req.Nombre}
	if err := uc.repo.Create(ctx, e); err != nil {
		uc.logger.Error("[ESTADO_ESTANQUE] Error al crear estado de estanque", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ESTANQUE] Estado de estanque creado", slog.String("id", e.ID.String()))
	return e, nil
}

func (uc *estadoEstanqueUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoEstanque, error) {
	uc.logger.Debug("[ESTADO_ESTANQUE] Buscando estado de estanque por ID", slog.String("id", id.String()))
	res, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("[ESTADO_ESTANQUE] Error al buscar estado de estanque por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if res == nil {
		uc.logger.Warn("[ESTADO_ESTANQUE] Estado de estanque no encontrado por ID", slog.String("id", id.String()))
		return nil, nil
	}
	uc.logger.Info("[ESTADO_ESTANQUE] Estado de estanque encontrado por ID", slog.String("id", id.String()))
	return res, nil
}

func (uc *estadoEstanqueUseCase) List(ctx context.Context) ([]model.EstadoEstanque, error) {
	uc.logger.Debug("[ESTADO_ESTANQUE] Listando estados de estanque")
	res, err := uc.repo.List(ctx)
	if err != nil {
		uc.logger.Error("[ESTADO_ESTANQUE] Error al listar estados de estanque", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ESTANQUE] Estados de estanque listados con éxito", slog.Int("cantidad", len(res)))
	return res, nil
}

func (uc *estadoEstanqueUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoEstanqueRequest) (*model.EstadoEstanque, error) {
	uc.logger.Debug("[ESTADO_ESTANQUE] Actualizando estado de estanque", slog.String("id", id.String()))
	updated, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("[ESTADO_ESTANQUE] Error al actualizar estado de estanque", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ESTANQUE] Estado de estanque actualizado", slog.String("id", id.String()))
	return updated, nil
}

func (uc *estadoEstanqueUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Debug("[ESTADO_ESTANQUE] Eliminando estado de estanque", slog.String("id", id.String()))
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("[ESTADO_ESTANQUE] Error al eliminar estado de estanque", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	uc.logger.Info("[ESTADO_ESTANQUE] Estado de estanque eliminado con éxito", slog.String("id", id.String()))
	return nil
}
