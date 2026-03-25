package handlers

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JobProvider interface {
	GetJob(id string) (*models.TranslationJob, error)
	ListJobs(limit, offset int) ([]models.TranslationJob, int, error)
}

type JobHandler struct {
	Service JobProvider
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.Service.GetJob(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Job not found")
		return
	}
	utils.JSON(w, http.StatusOK, "success", "Job retrieved", job)
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	offset := (page - 1) * limit

	jobs, total, err := h.Service.ListJobs(limit, offset)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve jobs")
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Jobs retrieved", map[string]interface{}{
		"jobs":  jobs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
