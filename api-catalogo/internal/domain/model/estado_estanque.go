package model

import "github.com/google/uuid"

// EstadoEstanque representa el estado operativo de un estanque
type EstadoEstanque struct {
	ID     uuid.UUID `json:"id"`
	Codigo string    `json:"codigo,omitempty"`
	Nombre string    `json:"nombre,omitempty"`
}

// CreateEstadoEstanqueRequest es el payload para crear un nuevo estado
type CreateEstadoEstanqueRequest struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

// UpdateEstadoEstanqueRequest es el payload para actualizar un estado
type UpdateEstadoEstanqueRequest struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}
