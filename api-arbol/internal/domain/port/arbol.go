package port

import (
	"context"

	"c4-arbol/internal/domain/model"
	"github.com/google/uuid"
)

// ArbolRepository define los métodos de persistencia para Árbol.
type ArbolRepository interface {
	Create(ctx context.Context, arbol *model.Arbol) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Arbol, error)
	List(ctx context.Context, filter model.ArbolFilter) ([]model.Arbol, error)
	Update(ctx context.Context, arbol *model.Arbol) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ArbolUseCase define el contrato para los casos de uso de Árbol.
type ArbolUseCase interface {
	Create(ctx context.Context, req model.CreateArbolRequest) (*model.Arbol, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Arbol, error)
	List(ctx context.Context, filter model.ArbolFilter) ([]model.Arbol, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateArbolRequest) (*model.Arbol, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// HistorialRepository define los métodos de persistencia para el Historial de Estados.
type HistorialRepository interface {
	Create(ctx context.Context, h *model.HistorialEstadoArbol) error
	ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.HistorialEstadoArbol, error)
}

// HistorialUseCase define el contrato de casos de uso para el Historial de Estados.
type HistorialUseCase interface {
	Create(ctx context.Context, arbolID uuid.UUID, req model.CreateHistorialRequest) (*model.HistorialEstadoArbol, error)
	ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.HistorialEstadoArbol, error)
}

// MedicionRepository define los métodos de persistencia para las Mediciones.
type MedicionRepository interface {
	Create(ctx context.Context, m *model.MedicionArbol) error
	ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.MedicionArbol, error)
}

// MedicionUseCase define el contrato de casos de uso para las Mediciones.
type MedicionUseCase interface {
	Create(ctx context.Context, arbolID uuid.UUID, req model.CreateMedicionRequest) (*model.MedicionArbol, error)
	ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.MedicionArbol, error)
}
