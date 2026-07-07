package http

import (
	"encoding/json"
	"net/http"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/http/response"

	"github.com/google/uuid"
)

type EstadoEstanqueHandler struct {
	uc port.EstadoEstanqueUseCase
}

func NewEstadoEstanqueHandler(uc port.EstadoEstanqueUseCase) *EstadoEstanqueHandler {
	return &EstadoEstanqueHandler{uc: uc}
}

// Create godoc
// @Summary      Crear estado de estanque
// @Description  Crea un nuevo estado de estanque en el catálogo
// @Tags         Estados de Estanque
// @Accept       json
// @Produce      json
// @Param        request body      model.CreateEstadoEstanqueRequest true "Datos para crear el estado de estanque"
// @Success      201  {object}  response.APIResponse{data=model.EstadoEstanque}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-estanque [post]
func (h *EstadoEstanqueHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateEstadoEstanqueRequest
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
// @Summary      Obtener estado de estanque por ID
// @Description  Obtiene los detalles de un estado de estanque por su UUID
// @Tags         Estados de Estanque
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del estado (UUID)"
// @Success      200  {object}  response.APIResponse{data=model.EstadoEstanque}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-estanque/{id} [get]
func (h *EstadoEstanqueHandler) GetByID() http.HandlerFunc {
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
// @Summary      Listar estados de estanque
// @Description  Obtiene la lista de todos los estados de estanque registrados
// @Tags         Estados de Estanque
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]model.EstadoEstanque}
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-estanque [get]
func (h *EstadoEstanqueHandler) List() http.HandlerFunc {
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
// @Summary      Actualizar estado de estanque
// @Description  Actualiza los datos de un estado de estanque existente
// @Tags         Estados de Estanque
// @Accept       json
// @Produce      json
// @Param        id      path      string                           true  "ID del estado (UUID)"
// @Param        request body      model.UpdateEstadoEstanqueRequest true  "Datos a actualizar"
// @Success      200  {object}  response.APIResponse{data=model.EstadoEstanque}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-estanque/{id} [put]
func (h *EstadoEstanqueHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		var req model.UpdateEstadoEstanqueRequest
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
// @Summary      Eliminar estado de estanque
// @Description  Elimina un estado de estanque del catálogo
// @Tags         Estados de Estanque
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del estado (UUID)"
// @Success      200  {object}  response.APIResponse{data=map[string]string}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/estados-estanque/{id} [delete]
func (h *EstadoEstanqueHandler) Delete() http.HandlerFunc {
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
