package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"c4-zona/internal/domain/model"
	"c4-zona/internal/domain/port"
	"github.com/google/uuid"
)

type zonaRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewZonaRepository crea una nueva instancia del repositorio de Zona con Postgres.
func NewZonaRepository(db *sql.DB, log *slog.Logger) port.ZonaRepository {
	return &zonaRepository{
		db:  db,
		log: log.With(slog.String("component", "zona_repository")),
	}
}

func (r *zonaRepository) Create(ctx context.Context, zona *model.Zona) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("zona_id", zona.ID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de zona")

	query := `
		INSERT INTO territorio.zona (id, ecoparque_id, nombre, descripcion, area_m2, fecha_creacion, fecha_actualizacion)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		zona.ID,
		zona.EcoparqueID,
		zona.Nombre,
		zona.Descripcion,
		zona.AreaM2,
		zona.FechaCreacion,
		zona.FechaActualizacion,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar zona en base de datos",
			slog.Any("error", err),
		)
		return err
	}

	log.DebugContext(ctx, "INSERT de zona completado")
	return nil
}

func (r *zonaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error) {
	log := r.log.With(slog.String("operation", "GetByID"), slog.String("zona_id", id.String()))
	log.DebugContext(ctx, "ejecutando SELECT zona por ID")

	query := `
		SELECT id, ecoparque_id, nombre, descripcion, area_m2, fecha_creacion, fecha_actualizacion
		FROM territorio.zona
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var zona model.Zona
	err := row.Scan(
		&zona.ID,
		&zona.EcoparqueID,
		&zona.Nombre,
		&zona.Descripcion,
		&zona.AreaM2,
		&zona.FechaCreacion,
		&zona.FechaActualizacion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WarnContext(ctx, "zona no encontrada en base de datos")
			return nil, errors.New("zona no encontrada")
		}
		log.ErrorContext(ctx, "error al consultar zona por ID",
			slog.Any("error", err),
		)
		return nil, err
	}

	log.DebugContext(ctx, "zona encontrada en base de datos")
	return &zona, nil
}

func (r *zonaRepository) List(ctx context.Context) ([]model.Zona, error) {
	log := r.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "ejecutando SELECT de todas las zonas")

	query := `
		SELECT id, ecoparque_id, nombre, descripcion, area_m2, fecha_creacion, fecha_actualizacion
		FROM territorio.zona
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.ErrorContext(ctx, "error al listar zonas en base de datos",
			slog.Any("error", err),
		)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var zonas []model.Zona
	for rows.Next() {
		var zona model.Zona
		if err := rows.Scan(
			&zona.ID,
			&zona.EcoparqueID,
			&zona.Nombre,
			&zona.Descripcion,
			&zona.AreaM2,
			&zona.FechaCreacion,
			&zona.FechaActualizacion,
		); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de zona",
				slog.Any("error", err),
			)
			return nil, err
		}
		zonas = append(zonas, zona)
	}

	log.DebugContext(ctx, "SELECT de zonas completado",
		slog.Int("count", len(zonas)),
	)
	return zonas, nil
}

func (r *zonaRepository) Update(ctx context.Context, zona *model.Zona) error {
	log := r.log.With(slog.String("operation", "Update"), slog.String("zona_id", zona.ID.String()))
	log.DebugContext(ctx, "ejecutando UPDATE de zona")

	query := `
		UPDATE territorio.zona
		SET ecoparque_id = $1, nombre = $2, descripcion = $3, area_m2 = $4, fecha_actualizacion = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		zona.EcoparqueID,
		zona.Nombre,
		zona.Descripcion,
		zona.AreaM2,
		zona.FechaActualizacion,
		zona.ID,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al actualizar zona en base de datos",
			slog.Any("error", err),
		)
		return err
	}

	log.DebugContext(ctx, "UPDATE de zona completado")
	return nil
}

func (r *zonaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.log.With(slog.String("operation", "Delete"), slog.String("zona_id", id.String()))
	log.DebugContext(ctx, "ejecutando DELETE de zona")

	query := `DELETE FROM territorio.zona WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.ErrorContext(ctx, "error al eliminar zona en base de datos",
			slog.Any("error", err),
		)
		return err
	}

	log.DebugContext(ctx, "DELETE de zona completado")
	return nil
}
