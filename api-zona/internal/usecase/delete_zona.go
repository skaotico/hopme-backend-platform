package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Delete elimina una Zona por su identificador único.
func (uc *zonaUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	log := uc.log.With(slog.String("operation", "Delete"), slog.String("zona_id", id.String()))
	log.InfoContext(ctx, "iniciando eliminación de zona")

	if err := uc.repo.Delete(ctx, id); err != nil {
		log.ErrorContext(ctx, "error al eliminar zona",
			slog.Any("error", err),
		)
		return err
	}

	log.InfoContext(ctx, "zona eliminada exitosamente")
	return nil
}
