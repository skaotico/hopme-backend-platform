package usecase

import (
	"context"
	"log/slog"
	"time"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"c4-parque/internal/infra/observability"
	"github.com/google/uuid"
)

type ecoparqueUseCase struct {
	repo port.EcoparqueRepository
}

func NewEcoparqueUseCase(repo port.EcoparqueRepository) port.EcoparqueUseCase {
	return &ecoparqueUseCase{
		repo: repo,
	}
}

func (uc *ecoparqueUseCase) Create(ctx context.Context, req model.CreateEcoparqueRequest) (*model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Iniciando caso de uso Create", slog.String("nombre", req.Nombre))

	ecoparque := &model.Ecoparque{
		ID:            uuid.New(),
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		Direccion:     req.Direccion,
		FechaCreacion: time.Now().UTC(),
	}

	log.Info("Guardando nuevo ecoparque en repositorio", slog.String("id", ecoparque.ID.String()))
	err := uc.repo.Create(ctx, ecoparque)
	if err != nil {
		log.Error("Error en repositorio al crear ecoparque", slog.Any("error", err))
		return nil, err
	}

	log.Debug("Caso de uso Create finalizado con éxito", slog.String("id", ecoparque.ID.String()))
	return ecoparque, nil
}

func (uc *ecoparqueUseCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Iniciando caso de uso GetByID", slog.String("id", id.String()))

	log.Info("Obteniendo ecoparque desde repositorio", slog.String("id", id.String()))
	ecoparque, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("Error en repositorio al obtener ecoparque", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}

	log.Debug("Caso de uso GetByID finalizado con éxito", slog.String("id", ecoparque.ID.String()))
	return ecoparque, nil
}

func (uc *ecoparqueUseCase) List(ctx context.Context) ([]model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Iniciando caso de uso List")

	log.Info("Listando ecoparques desde repositorio")
	parques, err := uc.repo.List(ctx)
	if err != nil {
		log.Error("Error en repositorio al listar ecoparques", slog.Any("error", err))
		return nil, err
	}

	log.Debug("Caso de uso List finalizado con éxito", slog.Int("total_parques", len(parques)))
	return parques, nil
}

func (uc *ecoparqueUseCase) Update(ctx context.Context, id uuid.UUID, req model.UpdateEcoparqueRequest) (*model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Iniciando caso de uso Update", slog.String("id", id.String()))

	log.Info("Validando existencia de ecoparque para actualización", slog.String("id", id.String()))
	ecoparque, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.Warn("No se pudo obtener el ecoparque para actualizar", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}

	if req.Nombre != nil {
		log.Debug("Actualizando campo nombre")
		ecoparque.Nombre = *req.Nombre
	}
	if req.Descripcion != nil {
		log.Debug("Actualizando campo descripcion")
		ecoparque.Descripcion = req.Descripcion
	}
	if req.Direccion != nil {
		log.Debug("Actualizando campo direccion")
		ecoparque.Direccion = req.Direccion
	}

	now := time.Now().UTC()
	ecoparque.FechaActualizacion = &now

	log.Info("Guardando actualización en repositorio", slog.String("id", ecoparque.ID.String()))
	err = uc.repo.Update(ctx, ecoparque)
	if err != nil {
		log.Error("Error en repositorio al actualizar ecoparque", slog.String("id", ecoparque.ID.String()), slog.Any("error", err))
		return nil, err
	}

	log.Debug("Caso de uso Update finalizado con éxito", slog.String("id", ecoparque.ID.String()))
	return ecoparque, nil
}

func (uc *ecoparqueUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	log := observability.FromContext(ctx)
	log.Debug("Iniciando caso de uso Delete", slog.String("id", id.String()))

	log.Info("Eliminando ecoparque en repositorio", slog.String("id", id.String()))
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		log.Error("Error en repositorio al eliminar ecoparque", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}

	log.Debug("Caso de uso Delete finalizado con éxito", slog.String("id", id.String()))
	return nil
}
