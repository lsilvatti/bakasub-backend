package services

import (
	"bakasub-backend/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LogService struct {
	DB *sql.DB
}

func NewLogService(db *sql.DB) *LogService {
	return &LogService{DB: db}
}

func (s *LogService) CreateLog(level, module, message string, metadata map[string]any) error {
	var metaJSON []byte
	var err error

	if metadata != nil {
		metaJSON, err = json.Marshal(metadata)
		if err != nil {
			metaJSON = nil
		}
	}

	_, err = s.DB.Exec("INSERT INTO logs (level, module, message, metadata) VALUES ($1, $2, $3, $4)", level, module, message, string(metaJSON))
	return err
}

func (s *LogService) GetLogs(limit, offset int, level, module string) ([]models.LogEntry, int, error) {
	var conditions []string
	var args []any
	argID := 1

	if level != "" {
		conditions = append(conditions, fmt.Sprintf("level = $%d", argID))
		args = append(args, level)
		argID++
	}
	if module != "" {
		conditions = append(conditions, fmt.Sprintf("module = $%d", argID))
		args = append(args, module)
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM logs %s", whereClause)
	err := s.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id, level, module, message, metadata, timestamp FROM logs %s ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", whereClause, argID, argID+1)
	args = append(args, limit, offset)

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.LogEntry
	for rows.Next() {
		var log models.LogEntry
		var metaString sql.NullString
		var timestamp time.Time

		if err := rows.Scan(&log.ID, &log.Level, &log.Module, &log.Message, &metaString, &timestamp); err != nil {
			return nil, 0, err
		}

		log.Timestamp = timestamp.Format(time.RFC3339)

		if metaString.Valid && metaString.String != "" {
			json.Unmarshal([]byte(metaString.String), &log.Metadata)
		}

		logs = append(logs, log)
	}

	if logs == nil {
		logs = []models.LogEntry{}
	}

	return logs, total, nil
}

func (s *LogService) StartLogCleanupTask() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	s.pruneOldLogs()

	for range ticker.C {
		s.pruneOldLogs()
	}
}

func (s *LogService) pruneOldLogs() {
	var retentionDays int
	err := s.DB.QueryRow("SELECT log_retention_days FROM user_config LIMIT 1").Scan(&retentionDays)
	if err != nil {
		retentionDays = 7
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	_, err = s.DB.Exec("DELETE FROM logs WHERE timestamp < $1", cutoffTime)
	if err != nil {
		_ = s.CreateLog("error", "log_service", "Failed to prune old logs", map[string]any{"error": err.Error()})
	}
}
