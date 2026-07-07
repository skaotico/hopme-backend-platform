package http

import (
	"encoding/json"
	"net/http"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/http/response"

	"github.com/google/uuid"
)

type EstadoAguaHandler struct {
	uc port.EstadoAguaUseCase
}

func NewEstadoAguaHandler(uc port.EstadoAguaUseCase) *EstadoAguaHandler {
	return &EstadoAguaHandler{uc: uc}
}

// Create godoc
// @Summary      Crear estado de agua
// @Description  Crea un nuevo estado de calidad de agua en el catálogo
// @Tags         Estados de Agua
// @Accept       json
// @Produce      json
// @Param        request body      model.CreateEstadoAguaRequest true "Datos para crear el estado de agua"
// @Success      201  {object}  response.APIResponse{data=model.EstadoAgua}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-agua [post]
func (h *EstadoAguaHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateEstadoAguaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		res, err := h.uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear estado", err.Error())
			return
		}

		response.Success(w, http.StatusCreated, res)
	}
}

// GetByID godoc
// @Summary      Obtener estado de agua por ID
// @Description  Obtiene los detalles de un estado de agua por su UUID
// @Tags         Estados de Agua
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del estado (UUID)"
// @Success      200  {object}  response.APIResponse{data=model.EstadoAgua}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-agua/{id} [get]
func (h *EstadoAguaHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		res, err := h.uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al buscar estado", err.Error())
			return
		}

		if res == nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Estado no encontrado", "")
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// List godoc
// @Summary      Listar estados de agua
// @Description  Obtiene la lista de todos los estados de agua registrados
// @Tags         Estados de Agua
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]model.EstadoAgua}
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-agua [get]
func (h *EstadoAguaHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.List(r.Context())
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar estados", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Update godoc
// @Summary      Actualizar estado de agua
// @Description  Actualiza los datos de un estado de agua existente
// @Tags         Estados de Agua
// @Accept       json
// @Produce      json
// @Param        id      path      string                           true  "ID del estado (UUID)"
// @Param        request body      model.UpdateEstadoAguaRequest true  "Datos a actualizar"
// @Success      200  {object}  response.APIResponse{data=model.EstadoAgua}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-agua/{id} [put]
func (h *EstadoAguaHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		var req model.UpdateEstadoAguaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		res, err := h.uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar estado", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Delete godoc
// @Summary      Eliminar estado de agua
// @Description  Elimina un estado de agua del catálogo
// @Tags         Estados de Agua
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del estado (UUID)"
// @Success      200  {object}  response.APIResponse{data=map[string]string}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-agua/{id} [delete]
func (h *EstadoAguaHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		if err := h.uc.Delete(r.Context(), id); err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar estado", err.Error())
			return
		}

		response.Success(w, http.StatusOK, map[string]string{"message": "Estado eliminado correctamente"})
	}
}
