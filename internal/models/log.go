package models

import "time"

type LogEntry struct {
	ID        int            `json:"id,omitempty"`
	Level     string         `json:"level"`
	EventType string         `json:"event_type,omitempty"`
	Module    string         `json:"module"`
	Message   string         `json:"message"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Timestamp string         `json:"timestamp"`
}

type SystemLog struct {
	ID        int       `json:"id"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
