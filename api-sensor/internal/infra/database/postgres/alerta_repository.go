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

type alertaRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewAlertaRepository crea una nueva instancia del repositorio de alertas.
func NewAlertaRepository(db *sql.DB, log *slog.Logger) port.AlertaRepository {
	return &alertaRepository{
		db:  db,
		log: log.With(slog.String("component", "alerta_repository")),
	}
}

func (r *alertaRepository) Create(ctx context.Context, a *model.Alerta) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("sensor_id", a.SensorID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de alerta")

	query := `
		INSERT INTO iot.alerta (id, sensor_id, tipo, valor_detectado, umbral, fecha_alerta, estado, observacion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.SensorID, a.Tipo, a.ValorDetectado, a.Umbral, a.FechaAlerta, a.Estado, a.Observacion,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar alerta en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de alerta completado")
	return nil
}

func (r *alertaRepository) List(ctx context.Context, sensorID *uuid.UUID, estado *string) ([]model.Alerta, error) {
	log := r.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "ejecutando SELECT de alertas con filtros")

	query := `
		SELECT id, sensor_id, tipo, valor_detectado, umbral, fecha_alerta, estado, observacion
		FROM iot.alerta
		WHERE 1 = 1
	`
	var args []interface{}
	argCount := 1

	if sensorID != nil {
		query += fmt.Sprintf(" AND sensor_id = $%d", argCount)
		args = append(args, *sensorID)
		argCount++
	}

	if estado != nil {
		query += fmt.Sprintf(" AND estado = $%d", argCount)
		args = append(args, *estado)
		argCount++
	}

	query += " ORDER BY fecha_alerta DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.ErrorContext(ctx, "error al listar alertas en base de datos", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var alertas []model.Alerta
	for rows.Next() {
		var a model.Alerta
		if err := rows.Scan(
			&a.ID, &a.SensorID, &a.Tipo, &a.ValorDetectado, &a.Umbral, &a.FechaAlerta, &a.Estado, &a.Observacion,
		); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de alerta", slog.Any("error", err))
			return nil, err
		}
		alertas = append(alertas, a)
	}

	log.DebugContext(ctx, "SELECT de alertas completado", slog.Int("count", len(alertas)))
	return alertas, nil
}

func (r *alertaRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Alerta, error) {
	log := r.log.With(slog.String("operation", "GetByID"), slog.String("alerta_id", id.String()))
	log.DebugContext(ctx, "ejecutando SELECT de alerta por ID")

	query := `
		SELECT id, sensor_id, tipo, valor_detectado, umbral, fecha_alerta, estado, observacion
		FROM iot.alerta
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var a model.Alerta
	err := row.Scan(
		&a.ID, &a.SensorID, &a.Tipo, &a.ValorDetectado, &a.Umbral, &a.FechaAlerta, &a.Estado, &a.Observacion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WarnContext(ctx, "alerta no encontrada en base de datos")
			return nil, errors.New("alerta no encontrada")
		}
		log.ErrorContext(ctx, "error al consultar alerta por ID", slog.Any("error", err))
		return nil, err
	}

	log.DebugContext(ctx, "alerta encontrada en base de datos")
	return &a, nil
}

func (r *alertaRepository) Update(ctx context.Context, a *model.Alerta) error {
	log := r.log.With(slog.String("operation", "Update"), slog.String("alerta_id", a.ID.String()))
	log.DebugContext(ctx, "ejecutando UPDATE de alerta")

	query := `
		UPDATE iot.alerta
		SET estado = $1, observacion = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, a.Estado, a.Observacion, a.ID)
	if err != nil {
		log.ErrorContext(ctx, "error al actualizar alerta en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "UPDATE de alerta completado")
	return nil
}
