package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"c4-catalogo/internal/domain/model"
	"c4-catalogo/internal/domain/port"
	"c4-catalogo/internal/infra/observability"

	"github.com/google/uuid"
)

type estadoAguaRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewEstadoAguaRepository crea una instancia del adaptador PostgreSQL para estados de agua
func NewEstadoAguaRepository(db *sql.DB, logger *slog.Logger) port.EstadoAguaRepository {
	return &estadoAguaRepository{
		db:     db,
		logger: observability.WithComponent(logger, "repository.postgres.estado_agua"),
	}
}

func (r *estadoAguaRepository) Create(ctx context.Context, e *model.EstadoAgua) error {
	r.logger.Debug("[DATABASE] Insertando estado_agua", slog.String("id", e.ID.String()))
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO catalogo.estado_agua (id, codigo, nombre) VALUES ($1, $2, $3)`,
		e.ID, e.Codigo, e.Nombre,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al insertar estado_agua", slog.String("id", e.ID.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_agua insertado con éxito", slog.String("id", e.ID.String()))
	return nil
}

func (r *estadoAguaRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoAgua, error) {
	r.logger.Debug("[DATABASE] Buscando estado_agua por ID", slog.String("id", id.String()))
	e := &model.EstadoAgua{}
	var codigo, nombre sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre FROM catalogo.estado_agua WHERE id = $1`, id,
	).Scan(&e.ID, &codigo, &nombre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("[DATABASE] Estado_agua no encontrado por ID", slog.String("id", id.String()))
			return nil, nil
		}
		r.logger.Error("[DATABASE] Error al buscar estado_agua por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	r.logger.Debug("[DATABASE] Estado_agua encontrado", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoAguaRepository) List(ctx context.Context) ([]model.EstadoAgua, error) {
	r.logger.Debug("[DATABASE] Consultando lista de estados de agua")
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre FROM catalogo.estado_agua ORDER BY nombre`,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al listar estados de agua", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []model.EstadoAgua
	for rows.Next() {
		var e model.EstadoAgua
		var codigo, nombre sql.NullString
		if err := rows.Scan(&e.ID, &codigo, &nombre); err != nil {
			r.logger.Error("[DATABASE] Error al escanear fila de estado_agua", slog.Any("error", err))
			return nil, err
		}
		if codigo.Valid {
			e.Codigo = codigo.String
		}
		if nombre.Valid {
			e.Nombre = nombre.String
		}
		result = append(result, e)
	}
	r.logger.Debug("[DATABASE] Lista de estados de agua obtenida", slog.Int("cantidad", len(result)))
	return result, nil
}

func (r *estadoAguaRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoAguaRequest) (*model.EstadoAgua, error) {
	r.logger.Debug("[DATABASE] Actualizando estado_agua", slog.String("id", id.String()))
	e := &model.EstadoAgua{}
	var codigo, nombre sql.NullString
	err := r.db.QueryRowContext(ctx,
		`UPDATE catalogo.estado_agua SET codigo=$1, nombre=$2 WHERE id=$3
		 RETURNING id, codigo, nombre`,
		req.Codigo, req.Nombre, id,
	).Scan(&e.ID, &codigo, &nombre)
	if err != nil {
		r.logger.Error("[DATABASE] Error al actualizar estado_agua", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	r.logger.Debug("[DATABASE] Estado_agua actualizado con éxito", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoAguaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("[DATABASE] Eliminando estado_agua", slog.String("id", id.String()))
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalogo.estado_agua WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("[DATABASE] Error al eliminar estado_agua", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_agua eliminado con éxito", slog.String("id", id.String()))
	return nil
}
