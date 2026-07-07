package http

import (
	"encoding/json"
	"net/http"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"c4-parque/internal/infra/http/response"
	"github.com/google/uuid"
)

// HealthHandler maneja el healthcheck
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Parque Operativo"})
	})
}

// CreateParqueHandler maneja la creación de ecoparques
func CreateParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateEcoparqueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		if req.Nombre == "" {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El nombre es requerido", "")
			return
		}

		parque, err := uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear ecoparque", "")
			return
		}

		response.Success(w, http.StatusCreated, parque)
	})
}

// GetParqueHandler maneja la obtención de un ecoparque
func GetParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		parque, err := uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Ecoparque no encontrado", "")
			return
		}

		response.Success(w, http.StatusOK, parque)
	})
}

// ListParquesHandler maneja el listado
func ListParquesHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parques, err := uc.List(r.Context())
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar ecoparques", "")
			return
		}

		response.Success(w, http.StatusOK, parques)
	})
}

// UpdateParqueHandler maneja la actualización
func UpdateParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		var req model.UpdateEcoparqueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		parque, err := uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar ecoparque", "")
			return
		}

		response.Success(w, http.StatusOK, parque)
	})
}

// DeleteParqueHandler maneja el borrado
func DeleteParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		err = uc.Delete(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar ecoparque", "")
			return
		}

		response.Success(w, http.StatusOK, nil)
	})
}
