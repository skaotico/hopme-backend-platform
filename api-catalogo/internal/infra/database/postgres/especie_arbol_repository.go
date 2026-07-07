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

type especieArbolRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewEspecieArbolRepository crea una instancia del adaptador PostgreSQL para especies
func NewEspecieArbolRepository(db *sql.DB, logger *slog.Logger) port.EspecieArbolRepository {
	return &especieArbolRepository{
		db:     db,
		logger: observability.WithComponent(logger, "repository.postgres.especie_arbol"),
	}
}

func (r *especieArbolRepository) Create(ctx context.Context, e *model.EspecieArbol) error {
	r.logger.Debug("[DATABASE] Insertando especie_arbol", slog.String("id", e.ID.String()))
	query := `
		INSERT INTO catalogo.especie_arbol (id, nombre_comun, nombre_cientifico, familia, descripcion, activo)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.NombreComun, e.NombreCientifico, e.Familia, e.Descripcion, e.Activo,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al insertar especie_arbol", slog.String("id", e.ID.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Especie_arbol insertada con éxito", slog.String("id", e.ID.String()))
	return nil
}

func (r *especieArbolRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.EspecieArbol, error) {
	r.logger.Debug("[DATABASE] Buscando especie_arbol por ID", slog.String("id", id.String()))
	query := `
		SELECT id, nombre_comun, nombre_cientifico, familia, descripcion, activo
		FROM catalogo.especie_arbol
		WHERE id = $1
	`
	e := &model.EspecieArbol{}
	var nombreComun, familia, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &nombreComun, &e.NombreCientifico, &familia, &descripcion, &e.Activo,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("[DATABASE] Especie_arbol no encontrada por ID", slog.String("id", id.String()))
			return nil, nil
		}
		r.logger.Error("[DATABASE] Error al buscar especie_arbol por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if nombreComun.Valid {
		e.NombreComun = nombreComun.String
	}
	if familia.Valid {
		e.Familia = familia.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Especie_arbol encontrada", slog.String("id", id.String()))
	return e, nil
}

func (r *especieArbolRepository) List(ctx context.Context) ([]model.EspecieArbol, error) {
	r.logger.Debug("[DATABASE] Consultando lista de especies")
	query := `
		SELECT id, nombre_comun, nombre_cientifico, familia, descripcion, activo
		FROM catalogo.especie_arbol
		ORDER BY nombre_cientifico
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("[DATABASE] Error al listar especies", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []model.EspecieArbol
	for rows.Next() {
		var e model.EspecieArbol
		var nombreComun, familia, descripcion sql.NullString
		if err := rows.Scan(&e.ID, &nombreComun, &e.NombreCientifico, &familia, &descripcion, &e.Activo); err != nil {
			r.logger.Error("[DATABASE] Error al escanear fila de especie_arbol", slog.Any("error", err))
			return nil, err
		}
		if nombreComun.Valid {
			e.NombreComun = nombreComun.String
		}
		if familia.Valid {
			e.Familia = familia.String
		}
		if descripcion.Valid {
			e.Descripcion = descripcion.String
		}
		result = append(result, e)
	}
	r.logger.Debug("[DATABASE] Lista de especies obtenida", slog.Int("cantidad", len(result)))
	return result, nil
}

func (r *especieArbolRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateEspecieArbolRequest) (*model.EspecieArbol, error) {
	r.logger.Debug("[DATABASE] Actualizando especie_arbol", slog.String("id", id.String()))
	query := `
		UPDATE catalogo.especie_arbol
		SET nombre_comun = $1, nombre_cientifico = $2, familia = $3, descripcion = $4, activo = $5
		WHERE id = $6
		RETURNING id, nombre_comun, nombre_cientifico, familia, descripcion, activo
	`
	e := &model.EspecieArbol{}
	var nombreComun, familia, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx, query,
		req.NombreComun, req.NombreCientifico, req.Familia, req.Descripcion, req.Activo, id,
	).Scan(&e.ID, &nombreComun, &e.NombreCientifico, &familia, &descripcion, &e.Activo)
	if err != nil {
		r.logger.Error("[DATABASE] Error al actualizar especie_arbol", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if nombreComun.Valid {
		e.NombreComun = nombreComun.String
	}
	if familia.Valid {
		e.Familia = familia.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Especie_arbol actualizada con éxito", slog.String("id", id.String()))
	return e, nil
}

func (r *especieArbolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("[DATABASE] Eliminando especie_arbol", slog.String("id", id.String()))
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalogo.especie_arbol WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("[DATABASE] Error al eliminar especie_arbol", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Especie_arbol eliminada con éxito", slog.String("id", id.String()))
	return nil
}
