package port

import (
	"context"

	"c4-catalogo/internal/domain/model"

	"github.com/google/uuid"
)

// EspecieArbolUseCase define el Puerto de Entrada para la gestión de especies
type EspecieArbolUseCase interface {
	Create(ctx context.Context, req model.CreateEspecieArbolRequest) (*model.EspecieArbol, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.EspecieArbol, error)
	List(ctx context.Context) ([]model.EspecieArbol, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEspecieArbolRequest) (*model.EspecieArbol, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoArbolUseCase define el Puerto de Entrada para la gestión de estados de árbol
type EstadoArbolUseCase interface {
	Create(ctx context.Context, req model.CreateEstadoArbolRequest) (*model.EstadoArbol, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoArbol, error)
	List(ctx context.Context) ([]model.EstadoArbol, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoArbolRequest) (*model.EstadoArbol, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoEstanqueUseCase define el Puerto de Entrada para la gestión de estados de estanque
type EstadoEstanqueUseCase interface {
	Create(ctx context.Context, req model.CreateEstadoEstanqueRequest) (*model.EstadoEstanque, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoEstanque, error)
	List(ctx context.Context) ([]model.EstadoEstanque, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoEstanqueRequest) (*model.EstadoEstanque, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoAguaUseCase define el Puerto de Entrada para la gestión de estados de agua
type EstadoAguaUseCase interface {
	Create(ctx context.Context, req model.CreateEstadoAguaRequest) (*model.EstadoAgua, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.EstadoAgua, error)
	List(ctx context.Context) ([]model.EstadoAgua, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoAguaRequest) (*model.EstadoAgua, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TipoSensorUseCase define el Puerto de Entrada para la gestión de tipos de sensor
type TipoSensorUseCase interface {
	Create(ctx context.Context, req model.CreateTipoSensorRequest) (*model.TipoSensor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TipoSensor, error)
	List(ctx context.Context) ([]model.TipoSensor, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateTipoSensorRequest) (*model.TipoSensor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
