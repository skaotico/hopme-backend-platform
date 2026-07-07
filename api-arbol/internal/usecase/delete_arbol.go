package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (uc *arbolUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	log := uc.log.With(slog.String("operation", "Delete"), slog.String("arbol_id", id.String()))
	log.InfoContext(ctx, "iniciando eliminacion de arbol")

	if err := uc.repo.Delete(ctx, id); err != nil {
		log.ErrorContext(ctx, "error al eliminar arbol", slog.Any("error", err))
		return err
	}

	log.InfoContext(ctx, "arbol eliminado exitosamente")
	return nil
}
