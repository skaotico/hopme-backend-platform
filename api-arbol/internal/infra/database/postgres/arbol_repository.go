package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"c4-arbol/internal/domain/model"
	"c4-arbol/internal/domain/port"
	"github.com/google/uuid"
)

type arbolRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewArbolRepository crea un repositorio para la gestión de árboles.
func NewArbolRepository(db *sql.DB, log *slog.Logger) port.ArbolRepository {
	return &arbolRepository{
		db:  db,
		log: log.With(slog.String("component", "arbol_repository")),
	}
}

func (r *arbolRepository) Create(ctx context.Context, a *model.Arbol) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("arbol_id", a.ID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de arbol")

	query := `
		INSERT INTO flora.arbol (
			id, zona_id, especie_id, estado_id, codigo, latitud, longitud, 
			edad_estimada_anios, altura_m, ancho_copa_m, diametro_tronco_cm, 
			fecha_plantacion, fecha_registro, observaciones, fecha_creacion, fecha_actualizacion
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.ZonaID, a.EspecieID, a.EstadoID, a.Codigo, a.Latitud, a.Longitud,
		a.EdadEstimadaAnios, a.AlturaM, a.AnchoCopaM, a.DiametroTroncoCm,
		a.FechaPlantacion, a.FechaRegistro, a.Observaciones, a.FechaCreacion, a.FechaActualizacion,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar arbol en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de arbol completado")
	return nil
}

func (r *arbolRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Arbol, error) {
	log := r.log.With(slog.String("operation", "GetByID"), slog.String("arbol_id", id.String()))
	log.DebugContext(ctx, "ejecutando SELECT arbol por ID")

	query := `
		SELECT 
			id, zona_id, especie_id, estado_id, codigo, latitud, longitud, 
			edad_estimada_anios, altura_m, ancho_copa_m, diametro_tronco_cm, 
			fecha_plantacion, fecha_registro, observaciones, fecha_creacion, fecha_actualizacion
		FROM flora.arbol
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var a model.Arbol
	err := row.Scan(
		&a.ID, &a.ZonaID, &a.EspecieID, &a.EstadoID, &a.Codigo, &a.Latitud, &a.Longitud,
		&a.EdadEstimadaAnios, &a.AlturaM, &a.AnchoCopaM, &a.DiametroTroncoCm,
		&a.FechaPlantacion, &a.FechaRegistro, &a.Observaciones, &a.FechaCreacion, &a.FechaActualizacion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WarnContext(ctx, "arbol no encontrado en base de datos")
			return nil, errors.New("arbol no encontrado")
		}
		log.ErrorContext(ctx, "error al consultar arbol por ID", slog.Any("error", err))
		return nil, err
	}

	log.DebugContext(ctx, "arbol encontrado en base de datos")
	return &a, nil
}

func (r *arbolRepository) List(ctx context.Context, filter model.ArbolFilter) ([]model.Arbol, error) {
	log := r.log.With(slog.String("operation", "List"))
	log.DebugContext(ctx, "ejecutando SELECT de arboles con filtros")

	query := `
		SELECT 
			id, zona_id, especie_id, estado_id, codigo, latitud, longitud, 
			edad_estimada_anios, altura_m, ancho_copa_m, diametro_tronco_cm, 
			fecha_plantacion, fecha_registro, observaciones, fecha_creacion, fecha_actualizacion
		FROM flora.arbol
	`
	var rows *sql.Rows
	var err error

	if filter.ZonaID != nil {
		query += " WHERE zona_id = $1"
		rows, err = r.db.QueryContext(ctx, query, *filter.ZonaID)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		log.ErrorContext(ctx, "error al listar arboles en base de datos", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var arboles []model.Arbol
	for rows.Next() {
		var a model.Arbol
		if err := rows.Scan(
			&a.ID, &a.ZonaID, &a.EspecieID, &a.EstadoID, &a.Codigo, &a.Latitud, &a.Longitud,
			&a.EdadEstimadaAnios, &a.AlturaM, &a.AnchoCopaM, &a.DiametroTroncoCm,
			&a.FechaPlantacion, &a.FechaRegistro, &a.Observaciones, &a.FechaCreacion, &a.FechaActualizacion,
		); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de arbol", slog.Any("error", err))
			return nil, err
		}
		arboles = append(arboles, a)
	}

	log.DebugContext(ctx, "SELECT de arboles completado", slog.Int("count", len(arboles)))
	return arboles, nil
}

func (r *arbolRepository) Update(ctx context.Context, a *model.Arbol) error {
	log := r.log.With(slog.String("operation", "Update"), slog.String("arbol_id", a.ID.String()))
	log.DebugContext(ctx, "ejecutando UPDATE de arbol")

	query := `
		UPDATE flora.arbol
		SET 
			zona_id = $1, especie_id = $2, estado_id = $3, codigo = $4, 
			latitud = $5, longitud = $6, edad_estimada_anios = $7, 
			altura_m = $8, ancho_copa_m = $9, diametro_tronco_cm = $10, 
			fecha_plantacion = $11, fecha_registro = $12, observaciones = $13, 
			fecha_actualizacion = $14
		WHERE id = $15
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ZonaID, a.EspecieID, a.EstadoID, a.Codigo, a.Latitud, a.Longitud,
		a.EdadEstimadaAnios, a.AlturaM, a.AnchoCopaM, a.DiametroTroncoCm,
		a.FechaPlantacion, a.FechaRegistro, a.Observaciones, a.FechaActualizacion,
		a.ID,
	)
	if err != nil {
		log.ErrorContext(ctx, "error al actualizar arbol en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "UPDATE de arbol completado")
	return nil
}

func (r *arbolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.log.With(slog.String("operation", "Delete"), slog.String("arbol_id", id.String()))
	log.DebugContext(ctx, "ejecutando DELETE de arbol")

	query := `DELETE FROM flora.arbol WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.ErrorContext(ctx, "error al eliminar arbol en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "DELETE de arbol completado")
	return nil
}

// Historial Repository implementation
type historialRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewHistorialRepository(db *sql.DB, log *slog.Logger) port.HistorialRepository {
	return &historialRepository{
		db:  db,
		log: log.With(slog.String("component", "historial_repository")),
	}
}

func (r *historialRepository) Create(ctx context.Context, h *model.HistorialEstadoArbol) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("arbol_id", h.ArbolID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de historial_estado_arbol")

	query := `
		INSERT INTO flora.historial_estado_arbol (id, arbol_id, estado_id, observacion, fecha_registro)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, h.ID, h.ArbolID, h.EstadoID, h.Observacion, h.FechaRegistro)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar historial en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de historial completado")
	return nil
}

func (r *historialRepository) ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.HistorialEstadoArbol, error) {
	log := r.log.With(slog.String("operation", "ListByArbolID"), slog.String("arbol_id", arbolID.String()))
	log.DebugContext(ctx, "ejecutando SELECT historial de estados por arbol_id")

	query := `
		SELECT id, arbol_id, estado_id, observacion, fecha_registro
		FROM flora.historial_estado_arbol
		WHERE arbol_id = $1
		ORDER BY fecha_registro DESC
	`
	rows, err := r.db.QueryContext(ctx, query, arbolID)
	if err != nil {
		log.ErrorContext(ctx, "error al consultar historial de estados", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var historial []model.HistorialEstadoArbol
	for rows.Next() {
		var h model.HistorialEstadoArbol
		if err := rows.Scan(&h.ID, &h.ArbolID, &h.EstadoID, &h.Observacion, &h.FechaRegistro); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de historial", slog.Any("error", err))
			return nil, err
		}
		historial = append(historial, h)
	}

	log.DebugContext(ctx, "SELECT historial completado", slog.Int("count", len(historial)))
	return historial, nil
}

// Medicion Repository implementation
type medicionRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewMedicionRepository(db *sql.DB, log *slog.Logger) port.MedicionRepository {
	return &medicionRepository{
		db:  db,
		log: log.With(slog.String("component", "medicion_repository")),
	}
}

func (r *medicionRepository) Create(ctx context.Context, m *model.MedicionArbol) error {
	log := r.log.With(slog.String("operation", "Create"), slog.String("arbol_id", m.ArbolID.String()))
	log.DebugContext(ctx, "ejecutando INSERT de medicion_arbol")

	query := `
		INSERT INTO flora.medicion_arbol (id, arbol_id, altura_m, ancho_copa_m, diametro_tronco_cm, observaciones, fecha_medicion)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, m.ID, m.ArbolID, m.AlturaM, m.AnchoCopaM, m.DiametroTroncoCm, m.Observaciones, m.FechaMedicion)
	if err != nil {
		log.ErrorContext(ctx, "error al insertar medicion en base de datos", slog.Any("error", err))
		return err
	}

	log.DebugContext(ctx, "INSERT de medicion completado")
	return nil
}

func (r *medicionRepository) ListByArbolID(ctx context.Context, arbolID uuid.UUID) ([]model.MedicionArbol, error) {
	log := r.log.With(slog.String("operation", "ListByArbolID"), slog.String("arbol_id", arbolID.String()))
	log.DebugContext(ctx, "ejecutando SELECT mediciones por arbol_id")

	query := `
		SELECT id, arbol_id, altura_m, ancho_copa_m, diametro_tronco_cm, observaciones, fecha_medicion
		FROM flora.medicion_arbol
		WHERE arbol_id = $1
		ORDER BY fecha_medicion DESC
	`
	rows, err := r.db.QueryContext(ctx, query, arbolID)
	if err != nil {
		log.ErrorContext(ctx, "error al consultar mediciones", slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var mediciones []model.MedicionArbol
	for rows.Next() {
		var m model.MedicionArbol
		if err := rows.Scan(&m.ID, &m.ArbolID, &m.AlturaM, &m.AnchoCopaM, &m.DiametroTroncoCm, &m.Observaciones, &m.FechaMedicion); err != nil {
			log.ErrorContext(ctx, "error al escanear fila de medicion", slog.Any("error", err))
			return nil, err
		}
		mediciones = append(mediciones, m)
	}

	log.DebugContext(ctx, "SELECT mediciones completado", slog.Int("count", len(mediciones)))
	return mediciones, nil
}
