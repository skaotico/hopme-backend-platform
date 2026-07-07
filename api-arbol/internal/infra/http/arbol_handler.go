package http

import (
	"encoding/json"
	"net/http"

	"c4-arbol/internal/domain/model"
	"c4-arbol/internal/domain/port"
	"c4-arbol/internal/infra/http/response"
	"github.com/google/uuid"
)

// HealthHandler maneja el healthcheck.
//
// @Summary     Health check del servicio
// @Description Verifica que el servicio Arbol está operativo
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Arbol Operativo"})
	})
}

// CreateArbolHandler maneja el registro de un árbol.
//
// @Summary     Registrar un árbol
// @Description Registra un nuevo árbol en una zona del ecoparque
// @Tags        arboles
// @Accept      json
// @Produce     json
// @Param       arbol body model.CreateArbolRequest true "Datos del árbol a registrar"
// @Success     201 {object} model.Arbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles [post]
func CreateArbolHandler(uc port.ArbolUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateArbolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		if req.ZonaID == uuid.Nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El campo zona_id es requerido", "")
			return
		}

		arbol, err := uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar árbol", "")
			return
		}

		response.Success(w, http.StatusCreated, arbol)
	})
}

// GetArbolHandler maneja la obtención de un árbol por ID.
//
// @Summary     Obtener un árbol por ID
// @Description Retorna los datos detallados de un árbol específico
// @Tags        arboles
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Success     200 {object} model.Arbol
// @Failure     400 {object} response.APIResponse
// @Failure     404 {object} response.APIResponse
// @Router      /arboles/{id} [get]
func GetArbolHandler(uc port.ArbolUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		arbol, err := uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Árbol no encontrado", "")
			return
		}

		response.Success(w, http.StatusOK, arbol)
	})
}

// ListArbolesHandler maneja el listado y filtrado de árboles.
//
// @Summary     Listar árboles
// @Description Retorna la lista de árboles filtrada opcionalmente por zona_id
// @Tags        arboles
// @Produce     json
// @Param       zona_id query string false "Filtrar por UUID de la zona"
// @Success     200 {array} model.Arbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles [get]
func ListArbolesHandler(uc port.ArbolUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filter model.ArbolFilter
		zonaIDStr := r.URL.Query().Get("zona_id")
		if zonaIDStr != "" {
			zonaID, err := uuid.Parse(zonaIDStr)
			if err != nil {
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "zona_id inválido", "")
				return
			}
			filter.ZonaID = &zonaID
		}

		arboles, err := uc.List(r.Context(), filter)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar árboles", "")
			return
		}

		response.Success(w, http.StatusOK, arboles)
	})
}

// UpdateArbolHandler maneja la actualización de un árbol.
//
// @Summary     Actualizar un árbol
// @Description Actualiza uno o más campos de un árbol existente
// @Tags        arboles
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Param       arbol body model.UpdateArbolRequest true "Campos a actualizar"
// @Success     200 {object} model.Arbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id} [put]
func UpdateArbolHandler(uc port.ArbolUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		var req model.UpdateArbolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		arbol, err := uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al actualizar árbol", "")
			return
		}

		response.Success(w, http.StatusOK, arbol)
	})
}

// DeleteArbolHandler maneja la eliminación de un árbol.
//
// @Summary     Eliminar un árbol
// @Description Elimina un árbol de la base de datos
// @Tags        arboles
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id} [delete]
func DeleteArbolHandler(uc port.ArbolUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		err = uc.Delete(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al eliminar árbol", "")
			return
		}

		response.Success(w, http.StatusOK, nil)
	})
}

// CreateHistorialHandler registra un cambio de estado en el historial.
//
// @Summary     Registrar cambio de estado
// @Description Registra un cambio de estado en el historial de un árbol y actualiza su estado actual
// @Tags        arboles
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Param       historial body model.CreateHistorialRequest true "Detalle del nuevo estado"
// @Success     201 {object} model.HistorialEstadoArbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id}/historial [post]
func CreateHistorialHandler(uc port.HistorialUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		arbolID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de árbol inválido", "")
			return
		}

		var req model.CreateHistorialRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		if req.EstadoID == uuid.Nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El campo estado_id es requerido", "")
			return
		}

		historial, err := uc.Create(r.Context(), arbolID, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar historial de estado", "")
			return
		}

		response.Success(w, http.StatusCreated, historial)
	})
}

// ListHistorialHandler lista el historial de estados de un árbol.
//
// @Summary     Listar historial de estados
// @Description Obtiene el historial completo de cambios de estado de un árbol específico
// @Tags        arboles
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Success     200 {array} model.HistorialEstadoArbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id}/historial [get]
func ListHistorialHandler(uc port.HistorialUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		arbolID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de árbol inválido", "")
			return
		}

		historial, err := uc.ListByArbolID(r.Context(), arbolID)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar historial de estados", "")
			return
		}

		response.Success(w, http.StatusOK, historial)
	})
}

// CreateMedicionHandler registra una medición dasométrica para un árbol.
//
// @Summary     Registrar medición dasométrica
// @Description Registra una nueva medición para un árbol y actualiza las dimensiones actuales del árbol
// @Tags        arboles
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Param       medicion body model.CreateMedicionRequest true "Detalle de la medición"
// @Success     201 {object} model.MedicionArbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id}/mediciones [post]
func CreateMedicionHandler(uc port.MedicionUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		arbolID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de árbol inválido", "")
			return
		}

		var req model.CreateMedicionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		medicion, err := uc.Create(r.Context(), arbolID, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar medición", "")
			return
		}

		response.Success(w, http.StatusCreated, medicion)
	})
}

// ListMedicionesHandler lista las mediciones de un árbol.
//
// @Summary     Listar mediciones de un árbol
// @Description Obtiene todas las mediciones registradas de un árbol específico
// @Tags        arboles
// @Produce     json
// @Param       id path string true "UUID del árbol"
// @Success     200 {array} model.MedicionArbol
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /arboles/{id}/mediciones [get]
func ListMedicionesHandler(uc port.MedicionUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		arbolID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de árbol inválido", "")
			return
		}

		mediciones, err := uc.ListByArbolID(r.Context(), arbolID)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar mediciones", "")
			return
		}

		response.Success(w, http.StatusOK, mediciones)
	})
}
