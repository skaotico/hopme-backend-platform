package model

import "github.com/google/uuid"

// EstadoArbol representa el estado de salud de un árbol
type EstadoArbol struct {
	ID          uuid.UUID `json:"id"`
	Codigo      string    `json:"codigo,omitempty"`
	Nombre      string    `json:"nombre,omitempty"`
	Descripcion string    `json:"descripcion,omitempty"`
}

// CreateEstadoArbolRequest es el payload para crear un nuevo estado
type CreateEstadoArbolRequest struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// UpdateEstadoArbolRequest es el payload para actualizar un estado
type UpdateEstadoArbolRequest struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}
