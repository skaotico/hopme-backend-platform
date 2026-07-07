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

type tipoSensorRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewTipoSensorRepository crea una instancia del adaptador PostgreSQL para tipos de sensor
func NewTipoSensorRepository(db *sql.DB, logger *slog.Logger) port.TipoSensorRepository {
	return &tipoSensorRepository{
		db:     db,
		logger: observability.WithComponent(logger, "repository.postgres.tipo_sensor"),
	}
}

func (r *tipoSensorRepository) Create(ctx context.Context, e *model.TipoSensor) error {
	r.logger.Debug("[DATABASE] Insertando tipo_sensor", slog.String("id", e.ID.String()))
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO catalogo.tipo_sensor (id, codigo, nombre, unidad_medida, descripcion) VALUES ($1, $2, $3, $4, $5)`,
		e.ID, e.Codigo, e.Nombre, e.UnidadMedida, e.Descripcion,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al insertar tipo_sensor", slog.String("id", e.ID.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Tipo_sensor insertado con éxito", slog.String("id", e.ID.String()))
	return nil
}

func (r *tipoSensorRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TipoSensor, error) {
	r.logger.Debug("[DATABASE] Buscando tipo_sensor por ID", slog.String("id", id.String()))
	e := &model.TipoSensor{}
	var codigo, nombre, unidadMedida, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, codigo, nombre, unidad_medida, descripcion FROM catalogo.tipo_sensor WHERE id = $1`, id,
	).Scan(&e.ID, &codigo, &nombre, &unidadMedida, &descripcion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("[DATABASE] Tipo_sensor no encontrado por ID", slog.String("id", id.String()))
			return nil, nil
		}
		r.logger.Error("[DATABASE] Error al buscar tipo_sensor por ID", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	if unidadMedida.Valid {
		e.UnidadMedida = unidadMedida.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Tipo_sensor encontrado", slog.String("id", id.String()))
	return e, nil
}

func (r *tipoSensorRepository) List(ctx context.Context) ([]model.TipoSensor, error) {
	r.logger.Debug("[DATABASE] Consultando lista de tipos de sensor")
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, codigo, nombre, unidad_medida, descripcion FROM catalogo.tipo_sensor ORDER BY nombre`,
	)
	if err != nil {
		r.logger.Error("[DATABASE] Error al listar tipos de sensor", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []model.TipoSensor
	for rows.Next() {
		var e model.TipoSensor
		var codigo, nombre, unidadMedida, descripcion sql.NullString
		if err := rows.Scan(&e.ID, &codigo, &nombre, &unidadMedida, &descripcion); err != nil {
			r.logger.Error("[DATABASE] Error al escanear fila de tipo_sensor", slog.Any("error", err))
			return nil, err
		}
		if codigo.Valid {
			e.Codigo = codigo.String
		}
		if nombre.Valid {
			e.Nombre = nombre.String
		}
		if unidadMedida.Valid {
			e.UnidadMedida = unidadMedida.String
		}
		if descripcion.Valid {
			e.Descripcion = descripcion.String
		}
		result = append(result, e)
	}
	r.logger.Debug("[DATABASE] Lista de tipos de sensor obtenida", slog.Int("cantidad", len(result)))
	return result, nil
}

func (r *tipoSensorRepository) Update(ctx context.Context, id uuid.UUID, req model.UpdateTipoSensorRequest) (*model.TipoSensor, error) {
	r.logger.Debug("[DATABASE] Actualizando tipo_sensor", slog.String("id", id.String()))
	e := &model.TipoSensor{}
	var codigo, nombre, unidadMedida, descripcion sql.NullString
	err := r.db.QueryRowContext(ctx,
		`UPDATE catalogo.tipo_sensor SET codigo=$1, nombre=$2, unidad_medida=$3, descripcion=$4 WHERE id=$5
		 RETURNING id, codigo, nombre, unidad_medida, descripcion`,
		req.Codigo, req.Nombre, req.UnidadMedida, req.Descripcion, id,
	).Scan(&e.ID, &codigo, &nombre, &unidadMedida, &descripcion)
	if err != nil {
		r.logger.Error("[DATABASE] Error al actualizar tipo_sensor", slog.String("id", id.String()), slog.Any("error", err))
		return nil, err
	}
	if codigo.Valid {
		e.Codigo = codigo.String
	}
	if nombre.Valid {
		e.Nombre = nombre.String
	}
	if unidadMedida.Valid {
		e.UnidadMedida = unidadMedida.String
	}
	if descripcion.Valid {
		e.Descripcion = descripcion.String
	}
	r.logger.Debug("[DATABASE] Tipo_sensor actualizado con éxito", slog.String("id", id.String()))
	return e, nil
}

func (r *tipoSensorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug("[DATABASE] Eliminando tipo_sensor", slog.String("id", id.String()))
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalogo.tipo_sensor WHERE id = $1`, id)
	if err != nil {
		r.logger.Error("[DATABASE] Error al eliminar tipo_sensor", slog.String("id", id.String()), slog.Any("error", err))
		return err
	}
	r.logger.Debug("[DATABASE] Tipo_sensor eliminado con éxito", slog.String("id", id.String()))
	return nil
}
