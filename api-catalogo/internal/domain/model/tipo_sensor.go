package model

import "github.com/google/uuid"

// TipoSensor representa un tipo de sensor del catálogo
type TipoSensor struct {
	ID           uuid.UUID `json:"id"`
	Codigo       string    `json:"codigo,omitempty"`
	Nombre       string    `json:"nombre,omitempty"`
	UnidadMedida string    `json:"unidad_medida,omitempty"`
	Descripcion  string    `json:"descripcion,omitempty"`
}

// CreateTipoSensorRequest es el payload para crear un nuevo tipo de sensor
type CreateTipoSensorRequest struct {
	Codigo       string `json:"codigo"`
	Nombre       string `json:"nombre"`
	UnidadMedida string `json:"unidad_medida"`
	Descripcion  string `json:"descripcion"`
}

// UpdateTipoSensorRequest es el payload para actualizar un tipo de sensor
type UpdateTipoSensorRequest struct {
	Codigo       string `json:"codigo"`
	Nombre       string `json:"nombre"`
	UnidadMedida string `json:"unidad_medida"`
	Descripcion  string `json:"descripcion"`
}
