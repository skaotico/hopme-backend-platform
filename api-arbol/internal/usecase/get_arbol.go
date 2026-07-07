package usecase

import (
	"context"
	"log/slog"

	"c4-arbol/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *arbolUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Arbol, error) {
	log := uc.log.With(slog.String("operation", "GetByID"), slog.String("arbol_id", id.String()))
	log.InfoContext(ctx, "buscando arbol por ID")

	arbol, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.WarnContext(ctx, "arbol no encontrado",
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "arbol encontrado")
	return arbol, nil
}
