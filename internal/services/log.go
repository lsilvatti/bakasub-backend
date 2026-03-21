package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"database/sql"
)

type LogService struct {
	DB *sql.DB
}

func NewLogService(db *sql.DB) *LogService {
	return &LogService{DB: db}
}

func (s *LogService) GetLogs(limit, offset int, level, module string) ([]models.SystemLog, int, error) {
	countQuery := "SELECT COUNT(*) FROM system_logs WHERE 1=1"
	dataQuery := "SELECT id, level, module, message, details, created_at FROM system_logs WHERE 1=1"

	args := []interface{}{}

	if level != "" {
		countQuery += " AND level = ?"
		dataQuery += " AND level = ?"
		args = append(args, level)
	}

	if module != "" {
		countQuery += " AND module = ?"
		dataQuery += " AND module = ?"
		args = append(args, module)
	}

	var total int
	err := s.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.LogError("logs", "Failed to count total logs for pagination", map[string]any{
			"error":  err.Error(),
			"level":  level,
			"module": module,
		})
		return nil, 0, err
	}

	dataQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.DB.Query(dataQuery, args...)
	if err != nil {
		utils.LogError("logs", "Failed to query system logs data", map[string]any{
			"error":  err.Error(),
			"limit":  limit,
			"offset": offset,
		})
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.SystemLog
	for rows.Next() {
		var l models.SystemLog
		var details sql.NullString

		if err := rows.Scan(&l.ID, &l.Level, &l.Module, &l.Message, &details, &l.CreatedAt); err != nil {
			utils.LogError("logs", "Failed to scan system log row", map[string]any{
				"error": err.Error(),
			})
			return nil, 0, err
		}

		if details.Valid {
			l.Details = details.String
		}

		logs = append(logs, l)
	}

	if logs == nil {
		logs = []models.SystemLog{}
	}

	return logs, total, nil
}
