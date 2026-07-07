package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"c4-parque/internal/infra/observability"
	"github.com/google/uuid"
)

type ecoparqueRepository struct {
	db *sql.DB
}

func NewEcoparqueRepository(db *sql.DB) port.EcoparqueRepository {
	return &ecoparqueRepository{db: db}
}

func (r *ecoparqueRepository) Create(ctx context.Context, ecoparque *model.Ecoparque) error {
	log := observability.FromContext(ctx)
	log.Debug("Ejecutando query INSERT en repositorio", slog.String("id", ecoparque.ID.String()))

	query := `
		INSERT INTO territorio.ecoparque (id, nombre, descripcion, direccion, fecha_creacion, fecha_actualizacion)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		ecoparque.ID,
		ecoparque.Nombre,
		ecoparque.Descripcion,
		ecoparque.Direccion,
		ecoparque.FechaCreacion,
		ecoparque.FechaActualizacion,
	)
	if err != nil {
		log.Error("Fallo al ejecutar query INSERT", slog.Any("error", err))
	} else {
		log.Debug("Query INSERT ejecutado correctamente")
	}
	return err
}

func (r *ecoparqueRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Ejecutando query SELECT por ID en repositorio", slog.String("id", id.String()))

	query := `
		SELECT id, nombre, descripcion, direccion, fecha_creacion, fecha_actualizacion
		FROM territorio.ecoparque
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var ecoparque model.Ecoparque
	err := row.Scan(
		&ecoparque.ID,
		&ecoparque.Nombre,
		&ecoparque.Descripcion,
		&ecoparque.Direccion,
		&ecoparque.FechaCreacion,
		&ecoparque.FechaActualizacion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("No se encontraron registros para el ID provisto", slog.String("id", id.String()))
			return nil, errors.New("ecoparque no encontrado")
		}
		log.Error("Fallo al ejecutar query SELECT por ID", slog.Any("error", err))
		return nil, err
	}
	log.Debug("Query SELECT por ID ejecutado correctamente", slog.String("id", id.String()))
	return &ecoparque, nil
}

func (r *ecoparqueRepository) List(ctx context.Context) ([]model.Ecoparque, error) {
	log := observability.FromContext(ctx)
	log.Debug("Ejecutando query SELECT para listar en repositorio")

	query := `
		SELECT id, nombre, descripcion, direccion, fecha_creacion, fecha_actualizacion
		FROM territorio.ecoparque
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Error("Fallo al ejecutar query SELECT (List)", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var parques []model.Ecoparque
	for rows.Next() {
		var ecoparque model.Ecoparque
		if err := rows.Scan(
			&ecoparque.ID,
			&ecoparque.Nombre,
			&ecoparque.Descripcion,
			&ecoparque.Direccion,
			&ecoparque.FechaCreacion,
			&ecoparque.FechaActualizacion,
		); err != nil {
			log.Error("Fallo al escanear fila de resultados", slog.Any("error", err))
			return nil, err
		}
		parques = append(parques, ecoparque)
	}
	log.Debug("Query SELECT (List) ejecutado correctamente", slog.Int("total_obtenidos", len(parques)))
	return parques, nil
}

func (r *ecoparqueRepository) Update(ctx context.Context, ecoparque *model.Ecoparque) error {
	log := observability.FromContext(ctx)
	log.Debug("Ejecutando query UPDATE en repositorio", slog.String("id", ecoparque.ID.String()))

	query := `
		UPDATE territorio.ecoparque
		SET nombre = $1, descripcion = $2, direccion = $3, fecha_actualizacion = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		ecoparque.Nombre,
		ecoparque.Descripcion,
		ecoparque.Direccion,
		ecoparque.FechaActualizacion,
		ecoparque.ID,
	)
	if err != nil {
		log.Error("Fallo al ejecutar query UPDATE", slog.String("id", ecoparque.ID.String()), slog.Any("error", err))
	} else {
		log.Debug("Query UPDATE ejecutado correctamente", slog.String("id", ecoparque.ID.String()))
	}
	return err
}

func (r *ecoparqueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := observability.FromContext(ctx)
	log.Debug("Ejecutando query DELETE en repositorio", slog.String("id", id.String()))

	query := `DELETE FROM territorio.ecoparque WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Error("Fallo al ejecutar query DELETE", slog.String("id", id.String()), slog.Any("error", err))
	} else {
		log.Debug("Query DELETE ejecutado correctamente", slog.String("id", id.String()))
	}
	return err
}
