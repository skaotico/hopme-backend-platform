package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-arbol/internal/domain/model"
	"github.com/google/uuid"
)

func (uc *arbolUseCase) Create(ctx context.Context, req model.CreateArbolRequest) (*model.Arbol, error) {
	log := uc.log.With(slog.String("operation", "Create"))
	log.InfoContext(ctx, "iniciando registro de arbol",
		slog.String("zona_id", req.ZonaID.String()),
	)

	now := time.Now().UTC()
	arbol := &model.Arbol{
		ID:                uuid.New(),
		ZonaID:            req.ZonaID,
		EspecieID:         req.EspecieID,
		EstadoID:          req.EstadoID,
		Codigo:            req.Codigo,
		Latitud:           req.Latitud,
		Longitud:          req.Longitud,
		EdadEstimadaAnios: req.EdadEstimadaAnios,
		AlturaM:           req.AlturaM,
		AnchoCopaM:        req.AnchoCopaM,
		DiametroTroncoCm:  req.DiametroTroncoCm,
		FechaPlantacion:   req.FechaPlantacion,
		FechaRegistro:     &now,
		Observaciones:     req.Observaciones,
		FechaCreacion:     now,
	}

	if err := uc.repo.Create(ctx, arbol); err != nil {
		log.ErrorContext(ctx, "error al crear arbol en repositorio",
			slog.String("arbol_id", arbol.ID.String()),
			slog.Any("error", err),
		)
		return nil, err
	}

	log.InfoContext(ctx, "arbol registrado exitosamente",
		slog.String("arbol_id", arbol.ID.String()),
	)
	return arbol, nil
}
