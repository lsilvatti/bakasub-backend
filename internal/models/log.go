package models

import "time"

type LogEntry struct {
	Level   string         `json:"level"`
	Module  string         `json:"module"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type SystemLog struct {
	ID        int       `json:"id"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
