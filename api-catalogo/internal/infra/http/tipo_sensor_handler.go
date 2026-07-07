package http

import (
	"encoding/json"
	"net/http"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/http/response"

	"github.com/google/uuid"
)

type TipoSensorHandler struct {
	uc port.TipoSensorUseCase
}

func NewTipoSensorHandler(uc port.TipoSensorUseCase) *TipoSensorHandler {
	return &TipoSensorHandler{uc: uc}
}

// Create godoc
// @Summary      Crear tipo de sensor
// @Description  Crea un nuevo tipo de sensor en el catálogo
// @Tags         Tipos de Sensor
// @Accept       json
// @Produce      json
// @Param        request body      model.CreateTipoSensorRequest true "Datos para crear el tipo de sensor"
// @Success      201  {object}  response.APIResponse{data=model.TipoSensor}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/tipos-sensor [post]
func (h *TipoSensorHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateTipoSensorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		res, err := h.uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear tipo de sensor", err.Error())
			return
		}

		response.Success(w, http.StatusCreated, res)
	}
}

// GetByID godoc
// @Summary      Obtener tipo de sensor por ID
// @Description  Obtiene los detalles de un tipo de sensor por su UUID
// @Tags         Tipos de Sensor
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del tipo de sensor (UUID)"
// @Success      200  {object}  response.APIResponse{data=model.TipoSensor}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/tipos-sensor/{id} [get]
func (h *TipoSensorHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		res, err := h.uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al buscar tipo de sensor", err.Error())
			return
		}

		if res == nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Tipo de sensor no encontrado", "")
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// List godoc
// @Summary      Listar tipos de sensor
// @Description  Obtiene la lista de todos los tipos de sensor registrados
// @Tags         Tipos de Sensor
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]model.TipoSensor}
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/tipos-sensor [get]
func (h *TipoSensorHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.List(r.Context())
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar tipos de sensor", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Update godoc
// @Summary      Actualizar tipo de sensor
// @Description  Actualiza los datos de un tipo de sensor existente
// @Tags         Tipos de Sensor
// @Accept       json
// @Produce      json
// @Param        id      path      string                           true  "ID del tipo de sensor (UUID)"
// @Param        request body      model.UpdateTipoSensorRequest true  "Datos a actualizar"
// @Success      200  {object}  response.APIResponse{data=model.TipoSensor}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/tipos-sensor/{id} [put]
func (h *TipoSensorHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		var req model.UpdateTipoSensorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		res, err := h.uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar tipo de sensor", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Delete godoc
// @Summary      Eliminar tipo de sensor
// @Description  Elimina un tipo de sensor del catálogo
// @Tags         Tipos de Sensor
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del tipo de sensor (UUID)"
// @Success      200  {object}  response.APIResponse{data=map[string]string}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/tipos-sensor/{id} [delete]
func (h *TipoSensorHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		if err := h.uc.Delete(r.Context(), id); err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar tipo de sensor", err.Error())
			return
		}

		response.Success(w, http.StatusOK, map[string]string{"message": "Tipo de sensor eliminado correctamente"})
	}
}
