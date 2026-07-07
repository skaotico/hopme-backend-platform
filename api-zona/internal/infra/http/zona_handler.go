package http

import (
	"encoding/json"
	"net/http"

	"c4-zona/internal/domain/model"
	"c4-zona/internal/domain/port"
	"c4-zona/internal/infra/http/response"
	"github.com/google/uuid"
)

// HealthHandler maneja el healthcheck
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Zona Operativo"})
	})
}

// CreateZonaHandler maneja la creación
func CreateZonaHandler(uc port.ZonaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateZonaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		if req.Nombre == "" || req.EcoparqueID == uuid.Nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El nombre y ecoparque_id son requeridos", "")
			return
		}

		zona, err := uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear zona", "")
			return
		}

		response.Success(w, http.StatusCreated, zona)
	})
}

// GetZonaHandler maneja la obtención
func GetZonaHandler(uc port.ZonaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		zona, err := uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Zona no encontrada", "")
			return
		}

		response.Success(w, http.StatusOK, zona)
	})
}

// ListZonasHandler maneja el listado
func ListZonasHandler(uc port.ZonaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zonas, err := uc.List(r.Context())
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar zonas", "")
			return
		}

		response.Success(w, http.StatusOK, zonas)
	})
}

// UpdateZonaHandler maneja la actualización
func UpdateZonaHandler(uc port.ZonaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		var req model.UpdateZonaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		zona, err := uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar zona", "")
			return
		}

		response.Success(w, http.StatusOK, zona)
	})
}

// DeleteZonaHandler maneja el borrado
func DeleteZonaHandler(uc port.ZonaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		err = uc.Delete(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar zona", "")
			return
		}

		response.Success(w, http.StatusOK, nil)
	})
}
