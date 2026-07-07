package model

import (
	"time"

	"github.com/google/uuid"
)

// Arbol representa un espécimen de árbol en el ecoparque.
type Arbol struct {
	ID                uuid.UUID  `json:"id"`
	ZonaID            uuid.UUID  `json:"zona_id"`
	EspecieID         *uuid.UUID `json:"especie_id,omitempty"`
	EstadoID          *uuid.UUID `json:"estado_id,omitempty"`
	Codigo            *string    `json:"codigo,omitempty"`
	Latitud           *float64   `json:"latitud,omitempty"`
	Longitud          *float64   `json:"longitud,omitempty"`
	EdadEstimadaAnios *int       `json:"edad_estimada_anios,omitempty"`
	AlturaM           *float64   `json:"altura_m,omitempty"`
	AnchoCopaM        *float64   `json:"ancho_copa_m,omitempty"`
	DiametroTroncoCm  *float64   `json:"diametro_tronco_cm,omitempty"`
	FechaPlantacion   *time.Time `json:"fecha_plantacion,omitempty"`
	FechaRegistro     *time.Time `json:"fecha_registro,omitempty"`
	Observaciones     *string    `json:"observaciones,omitempty"`
	FechaCreacion     time.Time  `json:"fecha_creacion"`
	FechaActualizacion *time.Time `json:"fecha_actualizacion,omitempty"`
}

// CreateArbolRequest define los datos necesarios para registrar un nuevo árbol.
type CreateArbolRequest struct {
	ZonaID            uuid.UUID  `json:"zona_id"`
	EspecieID         *uuid.UUID `json:"especie_id,omitempty"`
	EstadoID          *uuid.UUID `json:"estado_id,omitempty"`
	Codigo            *string    `json:"codigo,omitempty"`
	Latitud           *float64   `json:"latitud,omitempty"`
	Longitud          *float64   `json:"longitud,omitempty"`
	EdadEstimadaAnios *int       `json:"edad_estimada_anios,omitempty"`
	AlturaM           *float64   `json:"altura_m,omitempty"`
	AnchoCopaM        *float64   `json:"ancho_copa_m,omitempty"`
	DiametroTroncoCm  *float64   `json:"diametro_tronco_cm,omitempty"`
	FechaPlantacion   *time.Time `json:"fecha_plantacion,omitempty"`
	Observaciones     *string    `json:"observaciones,omitempty"`
}

// UpdateArbolRequest define los campos modificables de un árbol.
type UpdateArbolRequest struct {
	ZonaID            *uuid.UUID `json:"zona_id,omitempty"`
	EspecieID         *uuid.UUID `json:"especie_id,omitempty"`
	EstadoID          *uuid.UUID `json:"estado_id,omitempty"`
	Codigo            *string    `json:"codigo,omitempty"`
	Latitud           *float64   `json:"latitud,omitempty"`
	Longitud          *float64   `json:"longitud,omitempty"`
	EdadEstimadaAnios *int       `json:"edad_estimada_anios,omitempty"`
	AlturaM           *float64   `json:"altura_m,omitempty"`
	AnchoCopaM        *float64   `json:"ancho_copa_m,omitempty"`
	DiametroTroncoCm  *float64   `json:"diametro_tronco_cm,omitempty"`
	FechaPlantacion   *time.Time `json:"fecha_plantacion,omitempty"`
	Observaciones     *string    `json:"observaciones,omitempty"`
}

// HistorialEstadoArbol representa el registro de cambios de estado de un árbol.
type HistorialEstadoArbol struct {
	ID            uuid.UUID  `json:"id"`
	ArbolID       uuid.UUID  `json:"arbol_id"`
	EstadoID      uuid.UUID  `json:"estado_id"`
	Observacion   *string    `json:"observacion,omitempty"`
	FechaRegistro time.Time  `json:"fecha_registro"`
}

// CreateHistorialRequest define los datos para agregar un cambio de estado en el historial.
type CreateHistorialRequest struct {
	EstadoID    uuid.UUID `json:"estado_id"`
	Observacion *string   `json:"observacion,omitempty"`
}

// MedicionArbol representa los datos dasométricos recolectados de un árbol.
type MedicionArbol struct {
	ID               uuid.UUID `json:"id"`
	ArbolID          uuid.UUID `json:"arbol_id"`
	AlturaM          *float64  `json:"altura_m,omitempty"`
	AnchoCopaM       *float64  `json:"ancho_copa_m,omitempty"`
	DiametroTroncoCm *float64  `json:"diametro_tronco_cm,omitempty"`
	Observaciones    *string   `json:"observaciones,omitempty"`
	FechaMedicion    time.Time `json:"fecha_medicion"`
}

// CreateMedicionRequest define los datos para registrar una nueva medición.
type CreateMedicionRequest struct {
	AlturaM          *float64  `json:"altura_m,omitempty"`
	AnchoCopaM       *float64  `json:"ancho_copa_m,omitempty"`
	DiametroTroncoCm *float64  `json:"diametro_tronco_cm,omitempty"`
	Observaciones    *string   `json:"observaciones,omitempty"`
	FechaMedicion    *time.Time `json:"fecha_medicion,omitempty"`
}

// ArbolFilter define los filtros opcionales para listar árboles.
type ArbolFilter struct {
	ZonaID *uuid.UUID `json:"zona_id,omitempty"`
}

