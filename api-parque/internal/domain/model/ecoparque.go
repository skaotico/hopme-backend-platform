package model

import (
	"time"

	"github.com/google/uuid"
)

type Ecoparque struct {
	ID                 uuid.UUID  `json:"id"`
	Nombre             string     `json:"nombre"`
	Descripcion        *string    `json:"descripcion,omitempty"`
	Direccion          *string    `json:"direccion,omitempty"`
	FechaCreacion      time.Time  `json:"fecha_creacion"`
	FechaActualizacion *time.Time `json:"fecha_actualizacion,omitempty"`
}

type CreateEcoparqueRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
	Direccion   *string `json:"direccion,omitempty"`
}

type UpdateEcoparqueRequest struct {
	Nombre      *string `json:"nombre,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	Direccion   *string `json:"direccion,omitempty"`
}
