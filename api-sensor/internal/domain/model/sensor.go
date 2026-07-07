package model

import (
	"time"

	"github.com/google/uuid"
)

// Sensor representa un dispositivo físico de monitoreo en el ecoparque.
type Sensor struct {
	ID               uuid.UUID  `json:"id"`
	Codigo           *string    `json:"codigo,omitempty"`
	TipoSensorID     uuid.UUID  `json:"tipo_sensor_id"`
	ArbolID          *uuid.UUID `json:"arbol_id,omitempty"`
	EstanqueID       *uuid.UUID `json:"estanque_id,omitempty"`
	Fabricante       *string    `json:"fabricante,omitempty"`
	Modelo           *string    `json:"modelo,omitempty"`
	FechaInstalacion *time.Time `json:"fecha_instalacion,omitempty"`
	Activo           *bool      `json:"activo,omitempty"`
}

// CreateSensorRequest define los datos necesarios para registrar un nuevo sensor.
type CreateSensorRequest struct {
	Codigo           *string    `json:"codigo,omitempty"`
	TipoSensorID     uuid.UUID  `json:"tipo_sensor_id"`
	ArbolID          *uuid.UUID `json:"arbol_id,omitempty"`
	EstanqueID       *uuid.UUID `json:"estanque_id,omitempty"`
	Fabricante       *string    `json:"fabricante,omitempty"`
	Modelo           *string    `json:"modelo,omitempty"`
	FechaInstalacion *time.Time `json:"fecha_instalacion,omitempty"`
}

// UpdateSensorRequest define los campos modificables de un sensor.
type UpdateSensorRequest struct {
	Codigo           *string    `json:"codigo,omitempty"`
	TipoSensorID     *uuid.UUID `json:"tipo_sensor_id,omitempty"`
	ArbolID          *uuid.UUID `json:"arbol_id,omitempty"`
	EstanqueID       *uuid.UUID `json:"estanque_id,omitempty"`
	Fabricante       *string    `json:"fabricante,omitempty"`
	Modelo           *string    `json:"modelo,omitempty"`
	FechaInstalacion *time.Time `json:"fecha_instalacion,omitempty"`
	Activo           *bool      `json:"activo,omitempty"`
}

// SensorFilter define los filtros opcionales para listar sensores.
type SensorFilter struct {
	ArbolID    *uuid.UUID `json:"arbol_id,omitempty"`
	EstanqueID *uuid.UUID `json:"estanque_id,omitempty"`
	Activo     *bool      `json:"activo,omitempty"`
}

// LecturaSensor representa una lectura/métrica enviada por un sensor.
type LecturaSensor struct {
	ID                uuid.UUID `json:"id"`
	SensorID          uuid.UUID `json:"sensor_id"`
	Valor             float64   `json:"valor"`
	FechaLectura      time.Time `json:"fecha_lectura"`
	BateriaPorcentaje *int      `json:"bateria_porcentaje,omitempty"`
	Observacion       *string   `json:"observacion,omitempty"`
}

// CreateLecturaRequest define los datos para agregar una lectura de sensor.
type CreateLecturaRequest struct {
	Valor             float64   `json:"valor"`
	FechaLectura      *time.Time `json:"fecha_lectura,omitempty"`
	BateriaPorcentaje *int      `json:"bateria_porcentaje,omitempty"`
	Observacion       *string   `json:"observacion,omitempty"`
}

// Alerta representa una anomalía o evento de umbral detectado por un sensor.
type Alerta struct {
	ID             uuid.UUID  `json:"id"`
	SensorID       uuid.UUID  `json:"sensor_id"`
	Tipo           string     `json:"tipo"`
	ValorDetectado *float64   `json:"valor_detectado,omitempty"`
	Umbral         *float64   `json:"umbral,omitempty"`
	FechaAlerta    time.Time  `json:"fecha_alerta"`
	Estado         string     `json:"estado"`
	Observacion    *string    `json:"observacion,omitempty"`
}

// CreateAlertaRequest define los datos necesarios para registrar una alerta.
type CreateAlertaRequest struct {
	Tipo           string     `json:"tipo"`
	ValorDetectado *float64   `json:"valor_detectado,omitempty"`
	Umbral         *float64   `json:"umbral,omitempty"`
	FechaAlerta    *time.Time `json:"fecha_alerta,omitempty"`
	Observacion    *string    `json:"observacion,omitempty"`
}

// UpdateAlertaStatusRequest define el payload para actualizar el estado de una alerta.
type UpdateAlertaStatusRequest struct {
	Estado      string  `json:"estado"`
	Observacion *string `json:"observacion,omitempty"`
}
