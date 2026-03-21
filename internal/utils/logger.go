package utils

import (
	"bakasub-backend/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var (
	dbConn  *sql.DB
	logChan = make(chan models.LogEntry, 500)
)

func processLogs() {
	for entry := range logChan {
		if dbConn != nil {
			var detailsJSON string
			if entry.Details != nil {
				detailsBytes, err := json.Marshal(entry.Details)
				if err == nil {
					detailsJSON = string(detailsBytes)
				}
			}

			_, _ = dbConn.Exec("INSERT INTO system_logs (level, event_type, module, message, details) VALUES ($1, $2, $3, $4, $5)", entry.Level, entry.EventType, entry.Module, entry.Message, detailsJSON)
		}
	}
}

func InitLogger(db *sql.DB) {
	dbConn = db

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	go processLogs()
}

func LogInfo(eventType, module, message string, details map[string]any) {
	slog.Info(message, "module", module, "details", details)

	logChan <- models.LogEntry{
		Level:     "INFO",
		EventType: eventType,
		Module:    module,
		Message:   message,
		Details:   details,
	}
}

func LogError(module, message string, details map[string]any) {
	slog.Error(message, "module", module, "details", details)

	logChan <- models.LogEntry{
		Level:     "ERROR",
		EventType: "error",
		Module:    module,
		Message:   message,
		Details:   details,
	}
}

func AutoPruneLogs() {
	executarLimpeza := func() {
		if dbConn == nil {
			return
		}

		var days int
		err := dbConn.QueryRow("SELECT log_retention_days FROM user_config LIMIT 1").Scan(&days)
		if err != nil {
			days = 7
		}

		modifier := fmt.Sprintf("-%d days", days)

		_, err = dbConn.Exec(`DELETE FROM system_logs WHERE created_at < datetime('now', ?)`, modifier)
		if err == nil {
			LogInfo("system", "info", "Old logs cleanup completed", map[string]any{"retention_days": days})
		}
	}

	executarLimpeza()

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		executarLimpeza()
	}
}
