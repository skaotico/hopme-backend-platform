package usecase

import (
	"context"
	"time"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"github.com/google/uuid"
)

type ecoparqueUseCase struct {
	repo port.EcoparqueRepository
}

func NewEcoparqueUseCase(repo port.EcoparqueRepository) port.EcoparqueUseCase {
	return &ecoparqueUseCase{
		repo: repo,
	}
}

func (uc *ecoparqueUseCase) Create(ctx context.Context, req model.CreateEcoparqueRequest) (*model.Ecoparque, error) {
	ecoparque := &model.Ecoparque{
		ID:            uuid.New(),
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		Direccion:     req.Direccion,
		FechaCreacion: time.Now().UTC(),
	}

	err := uc.repo.Create(ctx, ecoparque)
	if err != nil {
		return nil, err
	}

	return ecoparque, nil
}

func (uc *ecoparqueUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ecoparqueUseCase) List(ctx context.Context) ([]model.Ecoparque, error) {
	return uc.repo.List(ctx)
}

func (uc *ecoparqueUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEcoparqueRequest) (*model.Ecoparque, error) {
	ecoparque, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Nombre != nil {
		ecoparque.Nombre = *req.Nombre
	}
	if req.Descripcion != nil {
		ecoparque.Descripcion = req.Descripcion
	}
	if req.Direccion != nil {
		ecoparque.Direccion = req.Direccion
	}

	now := time.Now().UTC()
	ecoparque.FechaActualizacion = &now

	err = uc.repo.Update(ctx, ecoparque)
	if err != nil {
		return nil, err
	}

	return ecoparque, nil
}

func (uc *ecoparqueUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
