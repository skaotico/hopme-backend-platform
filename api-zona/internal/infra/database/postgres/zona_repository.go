package postgres

import (
	"context"
	"database/sql"
	"errors"

	"c4-zona/internal/domain/model"
	"c4-zona/internal/domain/port"
	"github.com/google/uuid"
)

type zonaRepository struct {
	db *sql.DB
}

func NewZonaRepository(db *sql.DB) port.ZonaRepository {
	return &zonaRepository{db: db}
}

func (r *zonaRepository) Create(ctx context.Context, zona *model.Zona) error {
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
	return err
}

func (r *zonaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Zona, error) {
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
			return nil, errors.New("zona no encontrada")
		}
		return nil, err
	}
	return &zona, nil
}

func (r *zonaRepository) List(ctx context.Context) ([]model.Zona, error) {
	query := `
		SELECT id, ecoparque_id, nombre, descripcion, area_m2, fecha_creacion, fecha_actualizacion
		FROM territorio.zona
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			return nil, err
		}
		zonas = append(zonas, zona)
	}
	return zonas, nil
}

func (r *zonaRepository) Update(ctx context.Context, zona *model.Zona) error {
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
	return err
}

func (r *zonaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM territorio.zona WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
