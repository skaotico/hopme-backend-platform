package postgres

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations aplica las migraciones pendientes en la base de datos
func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations_catalogo",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	slog.Info("[DATABASE] Ejecutando migraciones pendientes...")
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("[DATABASE] Base de datos actualizada, sin cambios pendientes.")
			return nil
		}
		return err
	}

	slog.Info("[DATABASE] Migraciones ejecutadas satisfactoriamente.")
	return nil
}
