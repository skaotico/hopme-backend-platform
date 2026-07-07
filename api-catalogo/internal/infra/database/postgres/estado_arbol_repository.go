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

type estadoArbolRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewEstadoArbolRepository crea una instancia del adaptador PostgreSQL para estados de árbol
func NewEstadoArbolRepository(db *sql.DB, logger *slog.Logger) port.EstadoArbolRepository {
	return &estadoArbolRepository{
		db:     db,
		logger: observability.WithComponent(logger, "repository.postgres.estado_arbol"),
	}
}

func (r *estadoArbolRepository) Create(ctx context.Context, e *model.EstadoArbol) error {
	r.logger.Debug("[DATABASE] Insertando estado_arbol", slog.String("id", e.ID.String()))
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO catalogo.estado_arbol (id, codigo, nombre, descripcion) VALUES ($1, $2, $3, $4)`,
		e.ID, e.Codigo, e.Nombre, e.Descripcion,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al insertar estado_arbol", slog.String("id", e.ID.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_arbol insertado con éxito", slog.String("id", e.ID.String()))
	return nil
}

func (r *estadoArbolRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoArbol, error) {
	r.logger.Debug("[DATABASE] Buscando estado_arbol por ID", slog.String("id", id.String()))
	e := &model.EstadoArbol{}
	var codigo, nombre, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre, descripcion FROM catalogo.estado_arbol WHERE id = $1`, id,
	).Scan(&e.ID, &codigo, &nombre, &descripcion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("[DATABASE] Estado_arbol no encontrado por ID", slog.String("id", id.String()))
			return nil, nil
		}
		r.logger.Error("[DATABASE] Error al buscar estado_arbol por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Estado_arbol encontrado", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoArbolRepository) List(ctx context.Context) ([]model.EstadoArbol, error) {
	r.logger.Debug("[DATABASE] Consultando lista de estados de árbol")
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre, descripcion FROM catalogo.estado_arbol ORDER BY nombre`,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al listar estados de árbol", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []model.EstadoArbol
	for rows.Next() {
		var e model.EstadoArbol
		var codigo, nombre, descripcion sql.NullString
		if err := rows.Scan(&e.ID, &codigo, &nombre, &descripcion); err != nil {
			r.logger.Error("[DATABASE] Error al escanear fila de estado_arbol", slog.Any("error", err))
			return nil, err
		}
		if codigo.Valid {
			e.Codigo = codigo.String
		}
		if nombre.Valid {
			e.Nombre = nombre.String
		}
		if descripcion.Valid {
			e.Descripcion = descripcion.String
		}
		result = append(result, e)
	}
	r.logger.Debug("[DATABASE] Lista de estados de árbol obtenida", slog.Int("cantidad", len(result)))
	return result, nil
}

func (r *estadoArbolRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoArbolRequest) (*model.EstadoArbol, error) {
	r.logger.Debug("[DATABASE] Actualizando estado_arbol", slog.String("id", id.String()))
	e := &model.EstadoArbol{}
	var codigo, nombre, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx,
		`UPDATE catalogo.estado_arbol SET codigo=$1, nombre=$2, descripcion=$3 WHERE id=$4
		 RETURNING id, codigo, nombre, descripcion`,
		req.Codigo, req.Nombre, req.Descripcion, id,
	).Scan(&e.ID, &codigo, &nombre, &descripcion)
	if err != nil {
		r.logger.Error("[DATABASE] Error al actualizar estado_arbol", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Estado_arbol actualizado con éxito", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoArbolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("[DATABASE] Eliminando estado_arbol", slog.String("id", id.String()))
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalogo.estado_arbol WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("[DATABASE] Error al eliminar estado_arbol", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_arbol eliminado con éxito", slog.String("id", id.String()))
	return nil
}
