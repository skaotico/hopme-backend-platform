package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-arbol/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *arbolUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateArbolRequest) (*model.Arbol, error) {
	log := uc.log.With(slog.String("operation", "Update"), slog.String("arbol_id", id.String()))
	log.InfoContext(ctx, "iniciando actualizacion de arbol")

	arbol, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.WarnContext(ctx, "arbol no encontrado para actualizar", slog.Any("error", err))
		return nil, err
	}

	if req.ZonaID != nil {
		arbol.ZonaID = *req.ZonaID
	}
	if req.EspecieID != nil {
		arbol.EspecieID = req.EspecieID
	}
	if req.EstadoID != nil {
		arbol.EstadoID = req.EstadoID
	}
	if req.Codigo != nil {
		arbol.Codigo = req.Codigo
	}
	if req.Latitud != nil {
		arbol.Latitud = req.Latitud
	}
	if req.Longitud != nil {
		arbol.Longitud = req.Longitud
	}
	if req.EdadEstimadaAnios != nil {
		arbol.EdadEstimadaAnios = req.EdadEstimadaAnios
	}
	if req.AlturaM != nil {
		arbol.AlturaM = req.AlturaM
	}
	if req.AnchoCopaM != nil {
		arbol.AnchoCopaM = req.AnchoCopaM
	}
	if req.DiametroTroncoCm != nil {
		arbol.DiametroTroncoCm = req.DiametroTroncoCm
	}
	if req.FechaPlantacion != nil {
		arbol.FechaPlantacion = req.FechaPlantacion
	}
	if req.Observaciones != nil {
		arbol.Observaciones = req.Observaciones
	}

	now := time.Now().UTC()
	arbol.FechaActualizacion = &now

	if err := uc.repo.Update(ctx, arbol); err != nil {
		log.ErrorContext(ctx, "error al guardar actualizaciones del arbol", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "arbol actualizado exitosamente")
	return arbol, nil
}
