package postgres

import (
	"context"
	"database/sql"
	"errors"

	"c4-parque/internal/domain/model"
	"c4-parque/internal/domain/port"
	"github.com/google/uuid"
)

type ecoparqueRepository struct {
	db *sql.DB
}

func NewEcoparqueRepository(db *sql.DB) port.EcoparqueRepository {
	return &ecoparqueRepository{db: db}
}

func (r *ecoparqueRepository) Create(ctx context.Context, ecoparque *model.Ecoparque) error {
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
	return err
}

func (r *ecoparqueRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Ecoparque, error) {
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
			return nil, errors.New("ecoparque no encontrado")
		}
		return nil, err
	}
	return &ecoparque, nil
}

func (r *ecoparqueRepository) List(ctx context.Context) ([]model.Ecoparque, error) {
	query := `
		SELECT id, nombre, descripcion, direccion, fecha_creacion, fecha_actualizacion
		FROM territorio.ecoparque
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			return nil, err
		}
		parques = append(parques, ecoparque)
	}
	return parques, nil
}

func (r *ecoparqueRepository) Update(ctx context.Context, ecoparque *model.Ecoparque) error {
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
	return err
}

func (r *ecoparqueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM territorio.ecoparque WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
