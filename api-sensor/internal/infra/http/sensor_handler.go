package http

import (
	"encoding/json"
	"net/http"

	"c4-sensor/internal/domain/model"
	"c4-sensor/internal/domain/port"
	"c4-sensor/internal/infra/http/response"
	"github.com/google/uuid"
)

// HealthHandler maneja el healthcheck.
//
// @Summary     Health check del servicio
// @Description Verifica que el servicio Sensor está operativo
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Sensor Operativo"})
	})
}

// CreateSensorHandler maneja el registro de un sensor.
//
// @Summary     Registrar un sensor
// @Description Registra un nuevo sensor de monitoreo
// @Tags        sensores
// @Accept      json
// @Produce     json
// @Param       sensor body model.CreateSensorRequest true "Datos del sensor a registrar"
// @Success     201 {object} model.Sensor
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores [post]
func CreateSensorHandler(uc port.SensorUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateSensorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		if req.TipoSensorID == uuid.Nil {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El campo tipo_sensor_id es requerido", "")
			return
		}

		sensor, err := uc.Create(r.Context(), req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar sensor", "")
			return
		}

		response.Success(w, http.StatusCreated, sensor)
	})
}

// GetSensorHandler maneja la obtención de un sensor por ID.
//
// @Summary     Obtener un sensor por ID
// @Description Retorna los datos detallados de un sensor específico
// @Tags        sensores
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Success     200 {object} model.Sensor
// @Failure     400 {object} response.APIResponse
// @Failure     404 {object} response.APIResponse
// @Router      /sensores/{id} [get]
func GetSensorHandler(uc port.SensorUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		sensor, err := uc.GetByID(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Sensor no encontrado", "")
			return
		}

		response.Success(w, http.StatusOK, sensor)
	})
}

// ListSensoresHandler maneja el listado y filtrado de sensores.
//
// @Summary     Listar sensores
// @Description Retorna la lista de sensores filtrada opcionalmente por arbol_id o estanque_id
// @Tags        sensores
// @Produce     json
// @Param       arbol_id query string false "Filtrar por UUID de árbol"
// @Param       estanque_id query string false "Filtrar por UUID de estanque"
// @Param       activo query boolean false "Filtrar por estado activo"
// @Success     200 {array} model.Sensor
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores [get]
func ListSensoresHandler(uc port.SensorUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filter model.SensorFilter

		arbolIDStr := r.URL.Query().Get("arbol_id")
		if arbolIDStr != "" {
			arbolID, err := uuid.Parse(arbolIDStr)
			if err != nil {
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "arbol_id inválido", "")
				return
			}
			filter.ArbolID = &arbolID
		}

		estanqueIDStr := r.URL.Query().Get("estanque_id")
		if estanqueIDStr != "" {
			estanqueID, err := uuid.Parse(estanqueIDStr)
			if err != nil {
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "estanque_id inválido", "")
				return
			}
			filter.EstanqueID = &estanqueID
		}

		activoStr := r.URL.Query().Get("activo")
		if activoStr != "" {
			var activo bool
			if activoStr == "true" {
				activo = true
			} else if activoStr == "false" {
				activo = false
			} else {
				response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "activo debe ser true o false", "")
				return
			}
			filter.Activo = &activo
		}

		sensores, err := uc.List(r.Context(), filter)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar sensores", "")
			return
		}

		response.Success(w, http.StatusOK, sensores)
	})
}

// UpdateSensorHandler maneja la actualización de un sensor.
//
// @Summary     Actualizar un sensor
// @Description Actualiza uno o más campos de un sensor existente
// @Tags        sensores
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Param       sensor body model.UpdateSensorRequest true "Campos a actualizar"
// @Success     200 {object} model.Sensor
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores/{id} [put]
func UpdateSensorHandler(uc port.SensorUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		var req model.UpdateSensorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		sensor, err := uc.Update(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al actualizar sensor", "")
			return
		}

		response.Success(w, http.StatusOK, sensor)
	})
}

// DeleteSensorHandler maneja la eliminación de un sensor.
//
// @Summary     Eliminar un sensor
// @Description Elimina un sensor de la base de datos
// @Tags        sensores
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores/{id} [delete]
func DeleteSensorHandler(uc port.SensorUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID inválido", "")
			return
		}

		err = uc.Delete(r.Context(), id)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al eliminar sensor", "")
			return
		}

		response.Success(w, http.StatusOK, nil)
	})
}

// CreateLecturaHandler registra una nueva lectura para un sensor.
//
// @Summary     Registrar lectura de sensor
// @Description Registra una nueva lectura/medición reportada por un sensor
// @Tags        lecturas
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Param       lectura body model.CreateLecturaRequest true "Datos de la lectura"
// @Success     201 {object} model.LecturaSensor
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores/{id}/lecturas [post]
func CreateLecturaHandler(uc port.LecturaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		sensorID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de sensor inválido", "")
			return
		}

		var req model.CreateLecturaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		lectura, err := uc.Create(r.Context(), sensorID, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar lectura", "")
			return
		}

		response.Success(w, http.StatusCreated, lectura)
	})
}

// ListLecturasHandler obtiene las lecturas registradas para un sensor.
//
// @Summary     Listar lecturas de un sensor
// @Description Retorna el histórico de lecturas reportadas por un sensor específico
// @Tags        lecturas
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Success     200 {array} model.LecturaSensor
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores/{id}/lecturas [get]
func ListLecturasHandler(uc port.LecturaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		sensorID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de sensor inválido", "")
			return
		}

		lecturas, err := uc.ListBySensorID(r.Context(), sensorID)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar lecturas", "")
			return
		}

		response.Success(w, http.StatusOK, lecturas)
	})
}

// CreateAlertaHandler registra una nueva alerta para un sensor.
//
// @Summary     Registrar alerta de sensor
// @Description Registra un evento de alerta disparado por un sensor
// @Tags        alertas
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID del sensor"
// @Param       alerta body model.CreateAlertaRequest true "Datos de la alerta"
// @Success     201 {object} model.Alerta
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /sensores/{id}/alertas [post]
func CreateAlertaHandler(uc port.AlertaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		sensorID, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de sensor inválido", "")
			return
		}

		var req model.CreateAlertaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		if req.Tipo == "" {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El campo tipo es requerido", "")
			return
		}

		alerta, err := uc.Create(r.Context(), sensorID, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al registrar alerta", "")
			return
		}

		response.Success(w, http.StatusCreated, alerta)
	})
}

// ListAlertasHandler obtiene las alertas globales o filtradas por sensor y/o estado.
//
// @Summary     Listar alertas
// @Description Retorna la lista de alertas del sistema, con filtros opcionales
// @Tags        alertas
// @Produce     json
// @Param       sensor_id query string false "Filtrar por UUID de sensor"
// @Param       estado query string false "Filtrar por estado (e.g. activa, resuelta)"
// @Success     200 {array} model.Alerta
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /alertas [get]
func ListAlertasHandler(uc port.AlertaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sensorIDPtr *uuid.UUID
		sensorIDStr := r.URL.Query().Get("sensor_id")
		if sensorIDStr != "" {
			sensorID, err := uuid.Parse(sensorIDStr)
			if err != nil {
				response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "sensor_id inválido", "")
				return
			}
			sensorIDPtr = &sensorID
		}

		var estadoPtr *string
		estadoStr := r.URL.Query().Get("estado")
		if estadoStr != "" {
			estadoPtr = &estadoStr
		}

		alertas, err := uc.List(r.Context(), sensorIDPtr, estadoPtr)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al listar alertas", "")
			return
		}

		response.Success(w, http.StatusOK, alertas)
	})
}

// UpdateAlertaStatusHandler actualiza el estado y resolución de una alerta.
//
// @Summary     Actualizar estado de alerta
// @Description Actualiza el estado (activa, resuelta, etc.) y observaciones de una alerta
// @Tags        alertas
// @Accept      json
// @Produce     json
// @Param       id path string true "UUID de la alerta"
// @Param       estado body model.UpdateAlertaStatusRequest true "Datos de actualización de la alerta"
// @Success     200 {object} model.Alerta
// @Failure     400 {object} response.APIResponse
// @Failure     500 {object} response.APIResponse
// @Router      /alertas/{id}/estado [put]
func UpdateAlertaStatusHandler(uc port.AlertaUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidID, "ID de alerta inválido", "")
			return
		}

		var req model.UpdateAlertaStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Failure(w, http.StatusBadRequest, response.ErrCodeInvalidJSON, "Petición inválida", "")
			return
		}

		if req.Estado == "" {
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El campo estado es requerido", "")
			return
		}

		alerta, err := uc.UpdateStatus(r.Context(), id, req)
		if err != nil {
			response.Failure(w, http.StatusInternalServerError, response.ErrCodeInternal, "Error al actualizar alerta", "")
			return
		}

		response.Success(w, http.StatusOK, alerta)
	})
}
