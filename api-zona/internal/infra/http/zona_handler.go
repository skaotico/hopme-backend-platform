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
//
// @Summary     Health check del servicio
// @Description Verifica que el servicio Zona está operativo
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Zona Operativo"})
	})
}

// CreateZonaHandler maneja la creación de una zona
//
// @Summary     Crear una zona
// @Description Crea una nueva zona asociada a un ecoparque
// @Tags        zonas
// @Accept      json
// @Produce     json
// @Param       zona body model.CreateZonaRequest true "Datos de la zona a crear"
// @Success     201 {object} model.Zona
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /zonas [post]
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

// GetZonaHandler maneja la obtención de una zona por ID
//
// @Summary     Obtener una zona por ID
// @Description Retorna los datos de una zona específica
// @Tags        zonas
// @Produce     json
// @Param       id path string true "UUID de la zona"
// @Success     200 {object} model.Zona
// @Failure     400 {object} response.APIResponse
// @Failure     404 {object} response.APIResponse
// @Router      /zonas/{id} [get]
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

// ListZonasHandler maneja el listado de todas las zonas
//
// @Summary     Listar zonas
// @Description Retorna la lista completa de zonas registradas
// @Tags        zonas
// @Produce     json
// @Success     200 {array} model.Zona
// @Failure     500 {object} response.APIResponse
// @Router      /zonas [get]
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

// UpdateZonaHandler maneja la actualización de una zona
//
// @Summary     Actualizar una zona
// @Description Actualiza los campos de una zona existente
// @Tags        zonas
// @Accept      json
// @Produce     json
// @Param       id   path string                  true "UUID de la zona"
// @Param       zona body model.UpdateZonaRequest true "Campos a actualizar"
// @Success     200 {object} model.Zona
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /zonas/{id} [put]
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

// DeleteZonaHandler maneja el borrado de una zona
//
// @Summary     Eliminar una zona
// @Description Elimina una zona por su ID
// @Tags        zonas
// @Produce     json
// @Param       id path string true "UUID de la zona"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /zonas/{id} [delete]
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
