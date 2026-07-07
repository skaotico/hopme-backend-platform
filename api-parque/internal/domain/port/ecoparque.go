package port

import (
	"context"

	"c4-parque/internal/domain/model"
	"github.com/google/uuid"
)

type EcoparqueRepository interface {
	Create(ctx context.Context, ecoparque *model.Ecoparque) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error)
	List(ctx context.Context) ([]model.Ecoparque, error)
	Update(ctx context.Context, ecoparque *model.Ecoparque) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EcoparqueUseCase interface {
	Create(ctx context.Context, req model.CreateEcoparqueRequest) (*model.Ecoparque, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error)
	List(ctx context.Context) ([]model.Ecoparque, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEcoparqueRequest) (*model.Ecoparque, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
