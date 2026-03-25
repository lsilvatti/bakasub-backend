package utils

import (
	"bakasub-backend/internal/models"
	"context"
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

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

type CustomHandler struct {
	h slog.Handler
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{h: h.h.WithAttrs(attrs)}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{h: h.h.WithGroup(name)}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	levelColor := colorReset

	switch r.Level {
	case slog.LevelInfo:
		levelColor = colorCyan
	case slog.LevelWarn:
		levelColor = colorYellow
	case slog.LevelError:
		levelColor = colorRed
	}

	timeStr := r.Time.Format("15:04:05")

	module := "system"
	var detailsStr string
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "module" {
			module = a.Value.String()
		} else if a.Key == "details" {
			detailsStr = fmt.Sprintf(" %s%v%s", colorGray, a.Value, colorReset)
		}
		return true
	})

	fmt.Printf("%s%s%s | %s%-5s%s | %s%-10s%s | %s%s%s%s\n",
		colorGray, timeStr, colorReset,
		levelColor, level, colorReset,
		colorYellow, module, colorReset,
		colorReset, r.Message, colorReset,
		detailsStr,
	)

	return nil
}

func processLogs() {
	for entry := range logChan {
		if dbConn != nil {
			if entry.Metadata == nil {
				entry.Metadata = make(map[string]any)
			}
			if entry.EventType != "" {
				entry.Metadata["event_type"] = entry.EventType
			}

			var metadataJSON string
			metadataBytes, err := json.Marshal(entry.Metadata)
			if err == nil {
				metadataJSON = string(metadataBytes)
			}

			_, err = dbConn.Exec("INSERT INTO logs (level, module, message, metadata) VALUES ($1, $2, $3, $4)",
				entry.Level, entry.Module, entry.Message, metadataJSON)
			if err != nil {
				fmt.Printf("Erro ao salvar log no banco de dados: %v\n", err)
			}
		}
	}
}

func InitLogger(db *sql.DB) {
	dbConn = db

	logger := slog.New(&CustomHandler{
		h: slog.NewJSONHandler(os.Stdout, nil),
	})
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
		Metadata:  details,
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
		Metadata:  details,
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

		cutoffTime := time.Now().AddDate(0, 0, -days)

		_, err = dbConn.Exec(`DELETE FROM logs WHERE timestamp < $1`, cutoffTime)
		if err == nil {
			LogInfo("system", "info", "Old logs cleanup completed", map[string]any{"retention_days": days})
		} else {
			LogError("system", "Failed to prune old logs", map[string]any{"error": err.Error()})
		}
	}

	cleanupOldLogs()

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		cleanupOldLogs()
	}
}
