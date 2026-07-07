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

type estadoEstanqueRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewEstadoEstanqueRepository crea una instancia del adaptador PostgreSQL para estados de estanque
func NewEstadoEstanqueRepository(db *sql.DB, logger *slog.Logger) port.EstadoEstanqueRepository {
	return &estadoEstanqueRepository{
		db:     db,
		logger: observability.WithComponent(logger, "repository.postgres.estado_estanque"),
	}
}

func (r *estadoEstanqueRepository) Create(ctx context.Context, e *model.EstadoEstanque) error {
	r.logger.Debug("[DATABASE] Insertando estado_estanque", slog.String("id", e.ID.String()))
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO catalogo.estado_estanque (id, codigo, nombre) VALUES ($1, $2, $3)`,
		e.ID, e.Codigo, e.Nombre,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al insertar estado_estanque", slog.String("id", e.ID.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_estanque insertado con éxito", slog.String("id", e.ID.String()))
	return nil
}

func (r *estadoEstanqueRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.EstadoEstanque, error) {
	r.logger.Debug("[DATABASE] Buscando estado_estanque por ID", slog.String("id", id.String()))
	e := &model.EstadoEstanque{}
	var codigo, nombre sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre FROM catalogo.estado_estanque WHERE id = $1`, id,
	).Scan(&e.ID, &codigo, &nombre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("[DATABASE] Estado_estanque no encontrado por ID", slog.String("id", id.String()))
			return nil, nil
		}
		r.logger.Error("[DATABASE] Error al buscar estado_estanque por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	r.logger.Debug("[DATABASE] Estado_estanque encontrado", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoEstanqueRepository) List(ctx context.Context) ([]model.EstadoEstanque, error) {
	r.logger.Debug("[DATABASE] Consultando lista de estados de estanque")
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre FROM catalogo.estado_estanque ORDER BY nombre`,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al listar estados de estanque", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []model.EstadoEstanque
	for rows.Next() {
		var e model.EstadoEstanque
		var codigo, nombre sql.NullString
		if err := rows.Scan(&e.ID, &codigo, &nombre); err != nil {
			r.logger.Error("[DATABASE] Error al escanear fila de estado_estanque", slog.Any("error", err))
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
	r.logger.Debug("[DATABASE] Lista de estados de estanque obtenida", slog.Int("cantidad", len(result)))
	return result, nil
}

func (r *estadoEstanqueRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateEstadoEstanqueRequest) (*model.EstadoEstanque, error) {
	r.logger.Debug("[DATABASE] Actualizando estado_estanque", slog.String("id", id.String()))
	e := &model.EstadoEstanque{}
	var codigo, nombre sql.NullString
	err := r.db.QueryRowContext(ctx,
		`UPDATE catalogo.estado_estanque SET codigo=$1, nombre=$2 WHERE id=$3
		 RETURNING id, codigo, nombre`,
		req.Codigo, req.Nombre, id,
	).Scan(&e.ID, &codigo, &nombre)
	if err != nil {
		r.logger.Error("[DATABASE] Error al actualizar estado_estanque", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	r.logger.Debug("[DATABASE] Estado_estanque actualizado con éxito", slog.String("id", id.String()))
	return e, nil
}

func (r *estadoEstanqueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("[DATABASE] Eliminando estado_estanque", slog.String("id", id.String()))
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalogo.estado_estanque WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("[DATABASE] Error al eliminar estado_estanque", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Estado_estanque eliminado con éxito", slog.String("id", id.String()))
	return nil
}
