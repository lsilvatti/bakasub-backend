package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func InitializePostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("falha ao definir dialeto do goose: %w", err)
	}

	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		return nil, fmt.Errorf("falha ao aplicar migrations: %w", err)
	}

	return db, nil
}
