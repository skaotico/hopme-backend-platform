package usecase

import (
	"context"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type especieArbolUseCase struct {
	repo   port.EspecieArbolRepository
	logger *slog.Logger
}

// NewEspecieArbolUseCase crea una instancia del caso de uso de especies
func NewEspecieArbolUseCase(repo port.EspecieArbolRepository, logger *slog.Logger) port.EspecieArbolUseCase {
	return &especieArbolUseCase{
		repo:   repo,
		logger: observability.WithComponent(logger, "usecase.especie_arbol"),
	}
}

func (uc *especieArbolUseCase) Create(ctx context.Context, req model.CreateEspecieArbolRequest) (*model.EspecieArbol, error) {
	uc.logger.Debug("[ESPECIE_ARBOL] Creando nueva especie", slog.String("nombre_cientifico", req.NombreCientifico))

	e := &model.EspecieArbol{
		ID:               uuid.New(),
		NombreComun:      req.NombreComun,
		NombreCientifico: req.NombreCientifico,
		Familia:          req.Familia,
		Descripcion:      req.Descripcion,
		Activo:           req.Activo,
	}

	if err := uc.repo.Create(ctx, e); err != nil {
		uc.logger.Error("[ESPECIE_ARBOL] Error al crear especie", slog.Any("error", err))
		return nil, err
	}

	uc.logger.Info("[ESPECIE_ARBOL] Especie creada", slog.String("id", e.ID.String()))
	return e, nil
}

func (uc *especieArbolUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.EspecieArbol, error) {
	uc.logger.Debug("[ESPECIE_ARBOL] Buscando especie por ID", slog.String("id", id.String()))
	res, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("[ESPECIE_ARBOL] Error al buscar especie por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if res == nil {
		uc.logger.Warn("[ESPECIE_ARBOL] Especie no encontrada por ID", slog.String("id", id.String()))
		return nil, nil
	}
	uc.logger.Info("[ESPECIE_ARBOL] Especie encontrada por ID", slog.String("id", id.String()))
	return res, nil
}

func (uc *especieArbolUseCase) List(ctx context.Context) ([]model.EspecieArbol, error) {
	uc.logger.Debug("[ESPECIE_ARBOL] Listando especies")
	res, err := uc.repo.List(ctx)
	if err != nil {
		uc.logger.Error("[ESPECIE_ARBOL] Error al listar especies", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESPECIE_ARBOL] Especies listadas con éxito", slog.Int("cantidad", len(res)))
	return res, nil
}

func (uc *especieArbolUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEspecieArbolRequest) (*model.EspecieArbol, error) {
	uc.logger.Debug("[ESPECIE_ARBOL] Actualizando especie", slog.String("id", id.String()))
	updated, err := uc.repo.Update(ctx, id, req)
	if err != nil {
		uc.logger.Error("[ESPECIE_ARBOL] Error al actualizar especie", slog.Any("error", err))
		return nil, err
	}
	uc.logger.Info("[ESPECIE_ARBOL] Especie actualizada", slog.String("id", id.String()))
	return updated, nil
}

func (uc *especieArbolUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Debug("[ESPECIE_ARBOL] Eliminando especie", slog.String("id", id.String()))
	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("[ESPECIE_ARBOL] Error al eliminar especie", slog.Any("error", err))
		return err
	}
	uc.logger.Info("[ESPECIE_ARBOL] Especie eliminada", slog.String("id", id.String()))
	return nil
}
