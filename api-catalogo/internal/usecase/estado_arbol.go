package usecase

import (
	"context"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type estadoArbolUseCase struct {
	repo   port.EstadoArbolRepository
	logger *slog.Logger
}

// NewEstadoArbolUseCase crea una instancia del caso de uso de estados de árbol
func NewEstadoArbolUseCase(repo port.EstadoArbolRepository, logger *slog.Logger) port.EstadoArbolUseCase {
	return &estadoArbolUseCase{
		repo:   repo,
		logger: observability.WithComponent(logger, "usecase.estado_arbol"),
	}
}

func (uc *estadoArbolUseCase) Create(ctx context.Context, req model.CreateEstadoArbolRequest) (*model.EstadoArbol, error) {
	uc.logger.Debug("[ESTADO_ARBOL] Creando nuevo estado", slog.String("codigo", req.Codigo))
	e := &model.EstadoArbol{ID: uuid.New(), Codigo: req.Codigo, Nombre: req.Nombre, Descripcion: req.Descripcion}
	if err := uc.repo.Create(ctx, e); err != nil {
		uc.logger.Error("[ESTADO_ARBOL] Error al crear estado de árbol", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ARBOL] Estado de árbol creado", slog.String("id", e.ID.String()))
	return e, nil
}

func (uc *estadoArbolUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoArbol, error) {
	uc.logger.Debug("[ESTADO_ARBOL] Buscando estado de árbol por ID", slog.String("id", id.String()))
	res, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("[ESTADO_ARBOL] Error al buscar estado de árbol por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if res == nil {
		uc.logger.Warn("[ESTADO_ARBOL] Estado de árbol no encontrado por ID", slog.String("id", id.String()))
		return nil, nil
	}
	uc.logger.Info("[ESTADO_ARBOL] Estado de árbol encontrado por ID", slog.String("id", id.String()))
	return res, nil
}

func (uc *estadoArbolUseCase) List(ctx context.Context) ([]model.EstadoArbol, error) {
	uc.logger.Debug("[ESTADO_ARBOL] Listando estados de árbol")
	res, err := uc.repo.List(ctx)
	if err != nil {
		uc.logger.Error("[ESTADO_ARBOL] Error al listar estados de árbol", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ARBOL] Estados de árbol listados con éxito", slog.Int("cantidad", len(res)))
	return res, nil
}

func (uc *estadoArbolUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoArbolRequest) (*model.EstadoArbol, error) {
	uc.logger.Debug("[ESTADO_ARBOL] Actualizando estado de árbol", slog.String("id", id.String()))
	updated, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("[ESTADO_ARBOL] Error al actualizar estado de árbol", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESTADO_ARBOL] Estado de árbol actualizado", slog.String("id", id.String()))
	return updated, nil
}

func (uc *estadoArbolUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Debug("[ESTADO_ARBOL] Eliminando estado de árbol", slog.String("id", id.String()))
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("[ESTADO_ARBOL] Error al eliminar estado de árbol", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	uc.logger.Info("[ESTADO_ARBOL] Estado de árbol eliminado con éxito", slog.String("id", id.String()))
	return nil
}
