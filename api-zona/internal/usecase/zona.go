package usecase

import (
	"context"
	"time"

	"c4-zona/internal/domain/model"
	"c4-zona/internal/domain/port"
	"github.com/google/uuid"
)

type zonaUseCase struct {
	repo port.ZonaRepository
}

func NewZonaUseCase(repo port.ZonaRepository) port.ZonaUseCase {
	return &zonaUseCase{
		repo: repo,
	}
}

func (uc *zonaUseCase) Create(ctx context.Context, req model.CreateZonaRequest) (*model.Zona, error) {
	zona := &model.Zona{
		ID:            uuid.New(),
		EcoparqueID:   req.EcoparqueID,
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		AreaM2:        req.AreaM2,
		FechaCreacion: time.Now().UTC(),
	}

	err := uc.repo.Create(ctx, zona)
	if err != nil {
		return nil, err
	}

	return zona, nil
}

func (uc *zonaUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *zonaUseCase) List(ctx context.Context) ([]model.Zona, error) {
	return uc.repo.List(ctx)
}

func (uc *zonaUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateZonaRequest) (*model.Zona, error) {
	zona, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.EcoparqueID != nil {
		zona.EcoparqueID = *req.EcoparqueID
	}
	if req.Nombre != nil {
		zona.Nombre = *req.Nombre
	}
	if req.Descripcion != nil {
		zona.Descripcion = req.Descripcion
	}
	if req.AreaM2 != nil {
		zona.AreaM2 = req.AreaM2
	}

	now := time.Now().UTC()
	zona.FechaActualizacion = &now

	err = uc.repo.Update(ctx, zona)
	if err != nil {
		return nil, err
	}

	return zona, nil
}

func (uc *zonaUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
