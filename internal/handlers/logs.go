package handlers

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"net/http"
	"strconv"
)

type LogProvider interface {
	GetLogs(limit, offset int, level, module string) ([]models.LogEntry, int, error)
}

type LogHandler struct {
	Service LogProvider
}

func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	level := r.URL.Query().Get("level")
	module := r.URL.Query().Get("module")

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	offset := (page - 1) * limit

	logs, total, err := h.Service.GetLogs(limit, offset, level, module)
	if err != nil {
		utils.LogError("log_handler", "Failed to retrieve logs via service", map[string]any{
			"page":   page,
			"level":  level,
			"module": module,
			"error":  err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve logs: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Logs retrieved successfully", map[string]interface{}{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
