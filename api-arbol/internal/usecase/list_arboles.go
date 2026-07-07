package usecase

import (
	"context"
	"log/slog"

	"c4-arbol/internal/domain/model"
)

func (uc *arbolUseCase) List(ctx context.Context, filter model.ArbolFilter) ([]model.Arbol, error) {
	log := uc.log.With(slog.String("operation", "List"))
	if filter.ZonaID != nil {
		log = log.With(slog.String("filter_zona_id", filter.ZonaID.String()))
	}
	log.InfoContext(ctx, "listando arboles")

	arboles, err := uc.repo.List(ctx, filter)
	if err != nil {
		log.ErrorContext(ctx, "error al listar arboles", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "arboles listados exitosamente", slog.Int("count", len(arboles)))
	return arboles, nil
}
