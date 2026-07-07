package model

import "github.com/google/uuid"

// EspecieArbol representa una especie de árbol del catálogo
type EspecieArbol struct {
	ID              uuid.UUID `json:"id"`
	NombreComun     string    `json:"nombre_comun,omitempty"`
	NombreCientifico string   `json:"nombre_cientifico"`
	Familia         string    `json:"familia,omitempty"`
	Descripcion     string    `json:"descripcion,omitempty"`
	Activo          bool      `json:"activo"`
}

// CreateEspecieArbolRequest es el payload para crear una nueva especie
type CreateEspecieArbolRequest struct {
	NombreComun     string `json:"nombre_comun"`
	NombreCientifico string `json:"nombre_cientifico"`
	Familia         string `json:"familia"`
	Descripcion     string `json:"descripcion"`
	Activo          bool   `json:"activo"`
}

// UpdateEspecieArbolRequest es el payload para actualizar una especie
type UpdateEspecieArbolRequest struct {
	NombreComun     string `json:"nombre_comun"`
	NombreCientifico string `json:"nombre_cientifico"`
	Familia         string `json:"familia"`
	Descripcion     string `json:"descripcion"`
	Activo          bool   `json:"activo"`
}
