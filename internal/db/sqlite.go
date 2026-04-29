package db

import (
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var sqliteMigrationsFS embed.FS

func InitializeSQLite(dsn string) (*sql.DB, error) {
	if err := ensureSQLiteDirectory(dsn); err != nil {
		return nil, fmt.Errorf("falha ao preparar diretório do banco de dados: %w", err)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco de dados: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	if err := configureSQLite(db); err != nil {
		return nil, err
	}

	goose.SetBaseFS(sqliteMigrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("falha ao definir dialeto do goose: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("falha ao aplicar migrations: %w", err)
	}

	return db, nil
}

func configureSQLite(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("falha ao aplicar pragma %q: %w", pragma, err)
		}
	}

	return nil
}

func ensureSQLiteDirectory(dsn string) error {
	path := sqlitePathFromDSN(dsn)
	if path == "" || path == ":memory:" {
		return nil
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("falha ao criar diretório %q: %w", dir, err)
	}

	return nil
}

func sqlitePathFromDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	if dsn == ":memory:" || strings.HasPrefix(dsn, "file::memory:") {
		return ":memory:"
	}

	if !strings.HasPrefix(dsn, "file:") {
		return dsn
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		return trimSQLiteURIPath(strings.TrimPrefix(dsn, "file:"))
	}

	if parsed.Path != "" {
		return trimSQLiteURIPath(parsed.Path)
	}

	return trimSQLiteURIPath(parsed.Opaque)
}

func trimSQLiteURIPath(raw string) string {
	trimmed := raw
	if idx := strings.Index(trimmed, "?"); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	return trimmed
}
