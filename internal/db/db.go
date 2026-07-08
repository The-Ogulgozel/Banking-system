package db

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	log.Println("connecting to database...")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init pgx: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db:%w", err)
	}

	log.Println("database connected successfully")

	return pool, nil
}

func RunMigrations(dsn, migrationsPath string) error {
	log.Println("running database migrations...")

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return fmt.Errorf("failed to init migrate: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("no new migrations to apply")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("migrations applied successfully")

	return nil
}
