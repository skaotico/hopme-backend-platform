package port

import (
	"context"

	"c4-zona/internal/domain/model"
	"github.com/google/uuid"
)

type ZonaRepository interface {
	Create(ctx context.Context, zona *model.Zona) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error)
	List(ctx context.Context) ([]model.Zona, error)
	Update(ctx context.Context, zona *model.Zona) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ZonaUseCase interface {
	Create(ctx context.Context, req model.CreateZonaRequest) (*model.Zona, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error)
	List(ctx context.Context) ([]model.Zona, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateZonaRequest) (*model.Zona, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
