package usecase

import (
	"context"
	"log/slog"

	"c4-zona/internal/domain/model"
)

// List retorna todas las Zonas registradas.
func (uc *zonaUseCase) List(ctx context.Context) ([]model.Zona, error) {
	log := uc.log.With(slog.String("operation", "List"))
	log.InfoContext(ctx, "listando todas las zonas")

	data, err := uc.repo.List(ctx)
	if err != nil {
		log.ErrorContext(ctx, "error al listar zonas",
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "zonas listadas exitosamente",
		slog.Int("count", len(data)),
	)
	return data, nil
}
