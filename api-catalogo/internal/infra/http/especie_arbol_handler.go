package http

import (
	"encoding/json"
	"net/http"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/http/response"

	"github.com/google/uuid"
)

// EspecieArbolHandler exposes endpoints for EspecieArbol
type EspecieArbolHandler struct {
	uc port.EspecieArbolUseCase
}

func NewEspecieArbolHandler(uc port.EspecieArbolUseCase) *EspecieArbolHandler {
	return &EspecieArbolHandler{uc: uc}
}

// Create godoc
// @Summary      Crear especie de árbol
// @Description  Crea una nueva especie de árbol en el catálogo
// @Tags         Especies de Árbol
// @Accept       json
// @Produce      json
// @Param        request body      model.CreateEspecieArbolRequest true "Datos para crear la especie"
// @Success      201  {object}  response.APIResponse{data=model.EspecieArbol}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/especies [post]
func (h *EspecieArbolHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateEspecieArbolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		if req.NombreCientifico == "" {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El nombre científico es requerido", "")
			return
		}

		res, err := h.uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear especie", err.Error())
			return
		}

		response.Success(w, http.StatusCreated, res)
	}
}

// GetByID godoc
// @Summary      Obtener especie de árbol por ID
// @Description  Obtiene los detalles de una especie de árbol por su UUID
// @Tags         Especies de Árbol
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la especie (UUID)"
// @Success      200  {object}  response.APIResponse{data=model.EspecieArbol}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/especies/{id} [get]
func (h *EspecieArbolHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		res, err := h.uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al buscar especie", err.Error())
			return
		}

		if res == nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Especie no encontrada", "")
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// List godoc
// @Summary      Listar especies de árbol
// @Description  Obtiene la lista de todas las especies de árbol registradas
// @Tags         Especies de Árbol
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]model.EspecieArbol}
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/especies [get]
func (h *EspecieArbolHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.uc.List(r.Context())
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar especies", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Update godoc
// @Summary      Actualizar especie de árbol
// @Description  Actualiza los datos de una especie de árbol existente
// @Tags         Especies de Árbol
// @Accept       json
// @Produce      json
// @Param        id      path      string                           true  "ID de la especie (UUID)"
// @Param        request body      model.UpdateEspecieArbolRequest true  "Datos a actualizar"
// @Success      200  {object}  response.APIResponse{data=model.EspecieArbol}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/especies/{id} [put]
func (h *EspecieArbolHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		var req model.UpdateEspecieArbolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "JSON inválido", err.Error())
			return
		}

		res, err := h.uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar especie", err.Error())
			return
		}

		response.Success(w, http.StatusOK, res)
	}
}

// Delete godoc
// @Summary      Eliminar especie de árbol
// @Description  Elimina una especie de árbol del catálogo
// @Tags         Especies de Árbol
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la especie (UUID)"
// @Success      200  {object}  response.APIResponse{data=map[string]string}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Router       /api/v1/catalogo/especies/{id} [delete]
func (h *EspecieArbolHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", err.Error())
			return
		}

		if err := h.uc.Delete(r.Context(), id); err != nil {
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar especie", err.Error())
			return
		}

		response.Success(w, http.StatusOK, map[string]string{"message": "Especie eliminada correctamente"})
	}
}
