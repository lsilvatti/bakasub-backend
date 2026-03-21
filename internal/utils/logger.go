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
			if entry.Details == nil {
				entry.Details = make(map[string]any)
			}
			entry.Details["event_type"] = entry.EventType

			var detailsJSON string
			detailsBytes, err := json.Marshal(entry.Details)
			if err == nil {
				detailsJSON = string(detailsBytes)
			}

			_, _ = dbConn.Exec("INSERT INTO system_logs (level, module, message, details) VALUES ($1, $2, $3, $4)",
				entry.Level, entry.Module, entry.Message, detailsJSON)
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

func LogInfo(module, eventType, message string, details map[string]any) {
	if details == nil {
		details = make(map[string]any)
	}
	details["event_type"] = eventType

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
	if details == nil {
		details = make(map[string]any)
	}
	details["event_type"] = "error"

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
	cleanupOldLogs := func() {
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

	cleanupOldLogs()

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		cleanupOldLogs()
	}
}
