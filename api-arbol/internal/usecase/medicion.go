package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-arbol/internal/domain/model"
	"c4-arbol/internal/domain/port"
	"github.com/google/uuid"
)

type medicionUseCase struct {
	repo      port.MedicionRepository
	arbolRepo port.ArbolRepository
	log       *slog.Logger
}

// NewMedicionUseCase crea una instancia de caso de uso para las mediciones de árboles.
func NewMedicionUseCase(repo port.MedicionRepository, arbolRepo port.ArbolRepository, log *slog.Logger) port.MedicionUseCase {
	return &medicionUseCase{
		repo:      repo,
		arbolRepo: arbolRepo,
		log:       log.With(slog.String("component", "medicion_usecase")),
	}
}

func (uc *medicionUseCase) Create(ctx context.Context, arbolID uuid.UUID, req model.CreateMedicionRequest) (*model.MedicionArbol, error) {
	log := uc.log.With(slog.String("operation", "Create"), slog.String("arbol_id", arbolID.String()))
	log.InfoContext(ctx, "iniciando registro de medicion de arbol")

	// Validar existencia del árbol
	arbol, err := uc.arbolRepo.GetByID(ctx, arbolID)
	if err != nil {
		log.WarnContext(ctx, "arbol no encontrado para registrar medicion", slog.Any("error", err))
		return nil, err
	}

	fechaMedicion := time.Now().UTC()
	if req.FechaMedicion != nil {
		fechaMedicion = req.FechaMedicion.UTC()
	}

	m := &model.MedicionArbol{
		ID:               uuid.New(),
		ArbolID:          arbolID,
		AlturaM:          req.AlturaM,
		AnchoCopaM:       req.AnchoCopaM,
		DiametroTroncoCm: req.DiametroTroncoCm,
		Observaciones:    req.Observaciones,
		FechaMedicion:    fechaMedicion,
	}

	// Registrar la medición
	if err := uc.repo.Create(ctx, m); err != nil {
		log.ErrorContext(ctx, "error al insertar medicion en base de datos", slog.Any("error", err))
		return nil, err
	}

	// Actualizar dimensiones actuales del árbol si vienen en la medición
	updated := false
	if req.AlturaM != nil {
		arbol.AlturaM = req.AlturaM
		updated = true
	}
	if req.AnchoCopaM != nil {
		arbol.AnchoCopaM = req.AnchoCopaM
		updated = true
	}
	if req.DiametroTroncoCm != nil {
		arbol.DiametroTroncoCm = req.DiametroTroncoCm
		updated = true
	}

	if updated {
		now := time.Now().UTC()
		arbol.FechaActualizacion = &now
		if err := uc.arbolRepo.Update(ctx, arbol); err != nil {
			log.ErrorContext(ctx, "error al actualizar dimensiones en el arbol", slog.Any("error", err))
		}
	}

	log.InfoContext(ctx, "medicion registrada exitosamente")
	return m, nil
}

func (uc *medicionUseCase) ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.MedicionArbol, error) {
	log := uc.log.With(slog.String("operation", "ListByArbolID"), slog.String("arbol_id", arbolID.String()))
	log.InfoContext(ctx, "obteniendo mediciones del arbol")

	// Validar existencia
	if _, err := uc.arbolRepo.GetByID(ctx, arbolID); err != nil {
		log.WarnContext(ctx, "arbol no encontrado al consultar mediciones", slog.Any("error", err))
		return nil, err
	}

	mediciones, err := uc.repo.ListByArbolID(ctx, arbolID)
	if err != nil {
		log.ErrorContext(ctx, "error al listar mediciones", slog.Any("error", err))
		return nil, err
	}

	log.InfoContext(ctx, "mediciones listadas exitosamente", slog.Int("count", len(mediciones)))
	return mediciones, nil
}
