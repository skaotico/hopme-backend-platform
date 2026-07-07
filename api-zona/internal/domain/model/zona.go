package model

import (
	"time"

	"github.com/google/uuid"
)

type Zona struct {
	ID                 uuid.UUID  `json:"id"`
	EcoparqueID        uuid.UUID  `json:"ecoparque_id"`
	Nombre             string     `json:"nombre"`
	Descripcion        *string    `json:"descripcion,omitempty"`
	AreaM2             *float64   `json:"area_m2,omitempty"`
	FechaCreacion      time.Time  `json:"fecha_creacion"`
	FechaActualizacion *time.Time `json:"fecha_actualizacion,omitempty"`
}

type CreateZonaRequest struct {
	EcoparqueID uuid.UUID `json:"ecoparque_id"`
	Nombre      string    `json:"nombre"`
	Descripcion *string   `json:"descripcion,omitempty"`
	AreaM2      *float64  `json:"area_m2,omitempty"`
}

type UpdateZonaRequest struct {
	EcoparqueID *uuid.UUID `json:"ecoparque_id,omitempty"`
	Nombre      *string    `json:"nombre,omitempty"`
	Descripcion *string    `json:"descripcion,omitempty"`
	AreaM2      *float64   `json:"area_m2,omitempty"`
}
