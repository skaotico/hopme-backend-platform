package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"c4-sensor/internal/domain/model"
	"c4-sensor/internal/domain/port"
	"github.com/google/uuid"
)

type sensorRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewSensorRepository crea una nueva instancia del repositorio de sensores.
func NewSensorRepository(db *sql.DB, log *slog.Logger) port.SensorRepository {
	return &sensorRepository{
		db:  db,
		log: log.With(slog.String("component", "sensor_repository")),
	}
}

func (r *sensorRepository) Create(ctx context.Context, s *model.Sensor) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("sensor_id", s.ID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de sensor")

	query := `
		INSERT INTO iot.sensor (
			id, codigo, tipo_sensor_id, arbol_id, estanque_id, fabricante, modelo, fecha_instalacion, activo
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.Codigo, s.TipoSensorID, s.ArbolID, s.EstanqueID, s.Fabricante, s.Modelo, s.FechaInstalacion, s.Activo,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar sensor en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de sensor completado")
	return nil
}

func (r *sensorRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Sensor, error) {
	log := r.log.With(slog.String("operation", "GetByID"), slog.String("sensor_id", id.String()))
	log.DebugContext(ctx, "ejecutando SELECT de sensor por ID")

	query := `
		SELECT id, codigo, tipo_sensor_id, arbol_id, estanque_id, fabricante, modelo, fecha_instalacion, activo
		FROM iot.sensor
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s model.Sensor
	err := row.Scan(
		&s.ID, &s.Codigo, &s.TipoSensorID, &s.ArbolID, &s.EstanqueID, &s.Fabricante, &s.Modelo, &s.FechaInstalacion, &s.Activo,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WarnContext(ctx, "sensor no encontrado en base de datos")
			return nil, errors.New("sensor no encontrado")
		}
		log.ErrorContext(ctx, "error al consultar sensor por ID", slog.Any("error", err))
		return nil, err
	}

	log.DebugContext(ctx, "sensor encontrado en base de datos")
	return &s, nil
}

func (r *sensorRepository) List(ctx context.Context, filter model.SensorFilter) ([]model.Sensor, error) {
	log := r.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "ejecutando SELECT de sensores con filtros")

	query := `
		SELECT id, codigo, tipo_sensor_id, arbol_id, estanque_id, fabricante, modelo, fecha_instalacion, activo
		FROM iot.sensor
		WHERE 1 = 1
	`
	var args []interface{}
	argCount := 1

	if filter.ArbolID != nil {
		query += fmt.Sprintf(" AND arbol_id = $%d", argCount)
		args = append(args, *filter.ArbolID)
		argCount++
	}

	if filter.EstanqueID != nil {
		query += fmt.Sprintf(" AND estanque_id = $%d", argCount)
		args = append(args, *filter.EstanqueID)
		argCount++
	}

	if filter.Activo != nil {
		query += fmt.Sprintf(" AND activo = $%d", argCount)
		args = append(args, *filter.Activo)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.ErrorContext(ctx, "error al listar sensores en base de datos", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var sensores []model.Sensor
	for rows.Next() {
		var s model.Sensor
		if err := rows.Scan(
			&s.ID, &s.Codigo, &s.TipoSensorID, &s.ArbolID, &s.EstanqueID, &s.Fabricante, &s.Modelo, &s.FechaInstalacion, &s.Activo,
		); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de sensor", slog.Any("error", err))
			return nil, err
		}
		sensores = append(sensores, s)
	}

	log.DebugContext(ctx, "SELECT de sensores completado", slog.Int("count", len(sensores)))
	return sensores, nil
}

func (r *sensorRepository) Update(ctx context.Context, s *model.Sensor) error {
	log := r.log.With(slog.String("operation", "Update"), slog.String("sensor_id", s.ID.String()))
	log.DebugContext(ctx, "ejecutando UPDATE de sensor")

	query := `
		UPDATE iot.sensor
		SET codigo = $1, tipo_sensor_id = $2, arbol_id = $3, estanque_id = $4, 
		    fabricante = $5, modelo = $6, fecha_instalacion = $7, activo = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		s.Codigo, s.TipoSensorID, s.ArbolID, s.EstanqueID, s.Fabricante, s.Modelo, s.FechaInstalacion, s.Activo, s.ID,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al actualizar sensor en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "UPDATE de sensor completado")
	return nil
}

func (r *sensorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.log.With(slog.String("operation", "Delete"), slog.String("sensor_id", id.String()))
	log.DebugContext(ctx, "ejecutando DELETE de sensor")

	query := `DELETE FROM iot.sensor WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.ErrorContext(ctx, "error al eliminar sensor en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "DELETE de sensor completado")
	return nil
}
