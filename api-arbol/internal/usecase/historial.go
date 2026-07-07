package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-arbol/internal/domain/model"
	"c4-arbol/internal/domain/port"
	"github.com/google/uuid"
)

type historialUseCase struct {
	repo      port.HistorialRepository
	arbolRepo port.ArbolRepository
	log       *slog.Logger
}

// NewHistorialUseCase crea una instancia de caso de uso para el historial de estados de árboles.
func NewHistorialUseCase(repo port.HistorialRepository, arbolRepo port.ArbolRepository, log *slog.Logger) port.HistorialUseCase {
	return &historialUseCase{
		repo:      repo,
		arbolRepo: arbolRepo,
		log:       log.With(slog.String("component", "historial_usecase")),
	}
}

func (uc *historialUseCase) Create(ctx context.Context, arbolID uuid.UUID, req model.CreateHistorialRequest) (*model.HistorialEstadoArbol, error) {
	log := uc.log.With(
		slog.String("operation", "Create"),
		slog.String("arbol_id", arbolID.String()),
		slog.String("estado_id", req.EstadoID.String()),
	)
	log.InfoContext(ctx, "iniciando registro de cambio de estado de arbol")

	// Validar que el árbol existe
	arbol, err := uc.arbolRepo.GetByID(ctx, arbolID)
	if err != nil {
		log.WarnContext(ctx, "arbol no encontrado para cambiar estado", slog.Any("error", err))
		return nil, err
	}

	h := &model.HistorialEstadoArbol{
		ID:            uuid.New(),
		ArbolID:       arbolID,
		EstadoID:      req.EstadoID,
		Observacion:   req.Observacion,
		FechaRegistro: time.Now().UTC(),
	}

	// Registrar historial
	if err := uc.repo.Create(ctx, h); err != nil {
		log.ErrorContext(ctx, "error al insertar en historial", slog.Any("error", err))
		return nil, err
	}

	// Actualizar estado actual del árbol
	arbol.EstadoID = &req.EstadoID
	now := time.Now().UTC()
	arbol.FechaActualizacion = &now
	if err := uc.arbolRepo.Update(ctx, arbol); err != nil {
		log.ErrorContext(ctx, "error al actualizar estado actual en el arbol", slog.Any("error", err))
		// No revertimos para evitar inconsistencias de fallos, pero logueamos críticamente
	}

	log.InfoContext(ctx, "cambio de estado registrado exitosamente")
	return h, nil
}

func (uc *historialUseCase) ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.HistorialEstadoArbol, error) {
	log := uc.log.With(slog.String("operation", "ListByArbolID"), slog.String("arbol_id", arbolID.String()))
	log.InfoContext(ctx, "obteniendo historial de estados del arbol")

	// Validar existencia
	if _, err := uc.arbolRepo.GetByID(ctx, arbolID); err != nil {
		log.WarnContext(ctx, "arbol no encontrado al consultar historial", slog.Any("error", err))
		return nil, err
	}

	historial, err := uc.repo.ListByArbolID(ctx, arbolID)
	if err != nil {
		log.ErrorContext(ctx, "error al listar historial", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "historial de estados listado exitosamente", slog.Int("count", len(historial)))
	return historial, nil
}
