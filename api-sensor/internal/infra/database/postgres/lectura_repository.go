package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"c4-sensor/internal/domain/model"
	"c4-sensor/internal/domain/port"
	"github.com/google/uuid"
)

type lecturaRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewLecturaRepository crea una nueva instancia del repositorio de lecturas.
func NewLecturaRepository(db *sql.DB, log *slog.Logger) port.LecturaRepository {
	return &lecturaRepository{
		db:  db,
		log: log.With(slog.String("component", "lectura_repository")),
	}
}

func (r *lecturaRepository) Create(ctx context.Context, l *model.LecturaSensor) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("sensor_id", l.SensorID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de lectura_sensor")

	query := `
		INSERT INTO iot.lectura_sensor (id, sensor_id, valor, fecha_lectura, bateria_porcentaje, observacion)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.SensorID, l.Valor, l.FechaLectura, l.BateriaPorcentaje, l.Observacion,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar lectura de sensor en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de lectura_sensor completado")
	return nil
}

func (r *lecturaRepository) ListBySensorID(ctx context.Context, sensorID uuid.UUID) ([]model.LecturaSensor, error) {
	log := r.log.With(slog.String("operation", "ListBySensorID"), slog.String("sensor_id", sensorID.String()))
	log.DebugContext(ctx, "ejecutando SELECT de lecturas por sensor_id")

	query := `
		SELECT id, sensor_id, valor, fecha_lectura, bateria_porcentaje, observacion
		FROM iot.lectura_sensor
		WHERE sensor_id = $1
		ORDER BY fecha_lectura DESC
	`
	rows, err := r.db.QueryContext(ctx, query, sensorID)
	if err != nil {
		log.ErrorContext(ctx, "error al consultar lecturas de sensor", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var lecturas []model.LecturaSensor
	for rows.Next() {
		var l model.LecturaSensor
		if err := rows.Scan(
			&l.ID, &l.SensorID, &l.Valor, &l.FechaLectura, &l.BateriaPorcentaje, &l.Observacion,
		); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de lectura", slog.Any("error", err))
			return nil, err
		}
		lecturas = append(lecturas, l)
	}

	log.DebugContext(ctx, "SELECT de lecturas completado", slog.Int("count", len(lecturas)))
	return lecturas, nil
}
