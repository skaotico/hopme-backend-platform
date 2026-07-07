package port

import (
	"context"

	"c4-catalogo/internal/domain/model"

	"github.com/google/uuid"
)

// EspecieArbolRepository define el Puerto de Salida para la persistencia de especies
type EspecieArbolRepository interface {
	Create(ctx context.Context, e *model.EspecieArbol) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.EspecieArbol, error)
	List(ctx context.Context) ([]model.EspecieArbol, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEspecieArbolRequest) (*model.EspecieArbol, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoArbolRepository define el Puerto de Salida para la persistencia de estados de árbol
type EstadoArbolRepository interface {
	Create(ctx context.Context, e *model.EstadoArbol) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoArbol, error)
	List(ctx context.Context) ([]model.EstadoArbol, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoArbolRequest) (*model.EstadoArbol, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoEstanqueRepository define el Puerto de Salida para la persistencia de estados de estanque
type EstadoEstanqueRepository interface {
	Create(ctx context.Context, e *model.EstadoEstanque) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoEstanque, error)
	List(ctx context.Context) ([]model.EstadoEstanque, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoEstanqueRequest) (*model.EstadoEstanque, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EstadoAguaRepository define el Puerto de Salida para la persistencia de estados de agua
type EstadoAguaRepository interface {
	Create(ctx context.Context, e *model.EstadoAgua) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoAgua, error)
	List(ctx context.Context) ([]model.EstadoAgua, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoAguaRequest) (*model.EstadoAgua, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TipoSensorRepository define el Puerto de Salida para la persistencia de tipos de sensor
type TipoSensorRepository interface {
	Create(ctx context.Context, e *model.TipoSensor) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.TipoSensor, error)
	List(ctx context.Context) ([]model.TipoSensor, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateTipoSensorRequest) (*model.TipoSensor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
