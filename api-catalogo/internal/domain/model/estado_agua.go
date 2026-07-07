package model

import "github.com/google/uuid"

// EstadoAgua representa el estado de calidad del agua
type EstadoAgua struct {
	ID     uuid.UUID `json:"id"`
	Codigo string    `json:"codigo,omitempty"`
	Nombre string    `json:"nombre,omitempty"`
}

// CreateEstadoAguaRequest es el payload para crear un nuevo estado
type CreateEstadoAguaRequest struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

// UpdateEstadoAguaRequest es el payload para actualizar un estado
type UpdateEstadoAguaRequest struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}
