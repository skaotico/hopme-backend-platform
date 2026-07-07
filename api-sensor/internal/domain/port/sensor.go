package port

import (
	"context"

	"c4-sensor/internal/domain/model"
	"github.com/google/uuid"
)

// SensorRepository define los métodos de persistencia para Sensores.
type SensorRepository interface {
	Create(ctx context.Context, sensor *model.Sensor) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Sensor, error)
	List(ctx context.Context, filter model.SensorFilter) ([]model.Sensor, error)
	Update(ctx context.Context, sensor *model.Sensor) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SensorUseCase define el contrato para los casos de uso de Sensores.
type SensorUseCase interface {
	Create(ctx context.Context, req model.CreateSensorRequest) (*model.Sensor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Sensor, error)
	List(ctx context.Context, filter model.SensorFilter) ([]model.Sensor, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateSensorRequest) (*model.Sensor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// LecturaRepository define los métodos de persistencia para las lecturas de los sensores.
type LecturaRepository interface {
	Create(ctx context.Context, lectura *model.LecturaSensor) error
	ListBySensorID(ctx context.Context, sensorID uuid.UUID) ([]model.LecturaSensor, error)
}

// LecturaUseCase define el contrato de casos de uso para las lecturas de sensores.
type LecturaUseCase interface {
	Create(ctx context.Context, sensorID uuid.UUID, req model.CreateLecturaRequest) (*model.LecturaSensor, error)
	ListBySensorID(ctx context.Context, sensorID uuid.UUID) ([]model.LecturaSensor, error)
}

// AlertaRepository define los métodos de persistencia para las alertas.
type AlertaRepository interface {
	Create(ctx context.Context, alerta *model.Alerta) error
	List(ctx context.Context, sensorID *uuid.UUID, estado *string) ([]model.Alerta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Alerta, error)
	Update(ctx context.Context, alerta *model.Alerta) error
}

// AlertaUseCase define el contrato de casos de uso para las alertas.
type AlertaUseCase interface {
	Create(ctx context.Context, sensorID uuid.UUID, req model.CreateAlertaRequest) (*model.Alerta, error)
	List(ctx context.Context, sensorID *uuid.UUID, estado *string) ([]model.Alerta, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req model.UpdateAlertaStatusRequest) (*model.Alerta, error)
}
