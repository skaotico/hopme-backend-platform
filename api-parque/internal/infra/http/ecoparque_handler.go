package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"c4-parque/internal/infra/http/response"
	"c4-parque/internal/infra/observability"
	"github.com/google/uuid"
)

// HealthHandler maneja el healthcheck
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		log.Debug("Healthcheck solicitado")
		response.Success(w, http.StatusOK, map[string]string{"message": "Servicio Parque Operativo"})
	})
}

// CreateParqueHandler maneja la creación de ecoparques
func CreateParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		log.Info("Procesando petición para crear ecoparque")

		var req model.CreateEcoparqueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error al decodificar petición", slog.Any("error", err))
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		if req.Nombre == "" {
			log.Warn("Petición rechazada: nombre vacío")
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "El nombre es requerido", "")
			return
		}

		log.Debug("Datos recibidos para crear ecoparque", slog.String("nombre", req.Nombre))

		parque, err := uc.Create(r.Context(), req)
		if err != nil {
			log.Error("Fallo interno al crear ecoparque", slog.Any("error", err))
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al crear ecoparque", "")
			return
		}

		log.Info("Ecoparque creado exitosamente", slog.String("id", parque.ID.String()))
		response.Success(w, http.StatusCreated, parque)
	})
}

// GetParqueHandler maneja la obtención de un ecoparque
func GetParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		idStr := r.PathValue("id")
		log.Info("Procesando petición para obtener ecoparque", slog.String("id", idStr))

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("ID inválido recibido", slog.String("id", idStr), slog.Any("error", err))
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		log.Debug("Buscando ecoparque por ID", slog.String("id", id.String()))
		parque, err := uc.GetByID(r.Context(), id)
		if err != nil {
			log.Warn("Ecoparque no encontrado", slog.String("id", id.String()), slog.Any("error", err))
			response.Failure(w, http.StatusNotFound, "NOT_FOUND", "Ecoparque no encontrado", "")
			return
		}

		log.Info("Ecoparque obtenido exitosamente", slog.String("id", parque.ID.String()))
		response.Success(w, http.StatusOK, parque)
	})
}

// ListParquesHandler maneja el listado
func ListParquesHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		log.Info("Procesando petición para listar ecoparques")

		log.Debug("Invocando caso de uso List")
		parques, err := uc.List(r.Context())
		if err != nil {
			log.Error("Error al listar ecoparques", slog.Any("error", err))
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al listar ecoparques", "")
			return
		}

		log.Info("Listado de ecoparques obtenido", slog.Int("count", len(parques)))
		response.Success(w, http.StatusOK, parques)
	})
}

// UpdateParqueHandler maneja la actualización
func UpdateParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		idStr := r.PathValue("id")
		log.Info("Procesando petición para actualizar ecoparque", slog.String("id", idStr))

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("ID inválido recibido", slog.String("id", idStr), slog.Any("error", err))
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		var req model.UpdateEcoparqueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error al decodificar petición", slog.Any("error", err))
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "Petición inválida", "")
			return
		}

		log.Debug("Invocando caso de uso Update", slog.String("id", id.String()))
		parque, err := uc.Update(r.Context(), id, req)
		if err != nil {
			log.Error("Error al actualizar ecoparque", slog.String("id", id.String()), slog.Any("error", err))
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al actualizar ecoparque", "")
			return
		}

		log.Info("Ecoparque actualizado exitosamente", slog.String("id", parque.ID.String()))
		response.Success(w, http.StatusOK, parque)
	})
}

// DeleteParqueHandler maneja el borrado
func DeleteParqueHandler(uc port.EcoparqueUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := observability.FromContext(r.Context())
		idStr := r.PathValue("id")
		log.Info("Procesando petición para eliminar ecoparque", slog.String("id", idStr))

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("ID inválido recibido", slog.String("id", idStr), slog.Any("error", err))
			response.Failure(w, http.StatusBadRequest, "BAD_REQUEST", "ID inválido", "")
			return
		}

		log.Debug("Invocando caso de uso Delete", slog.String("id", id.String()))
		err = uc.Delete(r.Context(), id)
		if err != nil {
			log.Error("Error al eliminar ecoparque", slog.String("id", id.String()), slog.Any("error", err))
			response.Failure(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error al eliminar ecoparque", "")
			return
		}

		log.Info("Ecoparque eliminado exitosamente", slog.String("id", id.String()))
		response.Success(w, http.StatusOK, nil)
	})
}
