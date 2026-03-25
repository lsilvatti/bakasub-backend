package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"database/sql"
)

type JobService struct {
	DB *sql.DB
}

func NewJobService(db *sql.DB) *JobService {
	return &JobService{DB: db}
}

func (s *JobService) CreateJob(id, filePath, targetLang, preset, model string) error {
	_, err := s.DB.Exec(`
		INSERT INTO translation_jobs (id, status, file_path, target_lang, preset, model) 
		VALUES ($1, 'pending', $2, $3, $4, $5)`,
		id, filePath, targetLang, preset, model,
	)
	if err != nil {
		utils.LogError("job_service", "Failed to create job", map[string]any{"id": id, "error": err.Error()})
	}
	return err
}

func (s *JobService) UpdateTotalLines(id string, totalLines int) error {
	_, err := s.DB.Exec(`UPDATE translation_jobs SET total_lines = $1, status = 'processing', updated_at = CURRENT_TIMESTAMP WHERE id = $2`, totalLines, id)
	return err
}

func (s *JobService) IncrementProgress(id string, lines, pTokens, cTokens int, cost float64) error {
	_, err := s.DB.Exec(`
		UPDATE translation_jobs 
		SET processed_lines = processed_lines + $1, 
		    prompt_tokens = prompt_tokens + $2, 
		    completion_tokens = completion_tokens + $3, 
		    cost_usd = cost_usd + $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`,
		lines, pTokens, cTokens, cost, id,
	)
	if err != nil {
		utils.LogError("job_service", "Failed to increment job progress", map[string]any{"id": id, "error": err.Error()})
	}
	return err
}

func (s *JobService) UpdateStatus(id, status, errorMsg string) error {
	_, err := s.DB.Exec(`UPDATE translation_jobs SET status = $1, error_message = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`, status, errorMsg, id)
	if err != nil {
		utils.LogError("job_service", "Failed to update job status", map[string]any{"id": id, "status": status, "error": err.Error()})
	}
	return err
}

func (s *JobService) GetJob(id string) (*models.TranslationJob, error) {
	var j models.TranslationJob
	var errMsg sql.NullString

	err := s.DB.QueryRow(`
		SELECT id, status, file_path, target_lang, preset, model, total_lines, processed_lines, prompt_tokens, completion_tokens, cost_usd, error_message, created_at, updated_at 
		FROM translation_jobs WHERE id = $1`, id).
		Scan(&j.ID, &j.Status, &j.FilePath, &j.TargetLang, &j.Preset, &j.Model, &j.TotalLines, &j.ProcessedLines, &j.PromptTokens, &j.CompletionTokens, &j.CostUSD, &errMsg, &j.CreatedAt, &j.UpdatedAt)

	if err != nil {
		return nil, err
	}
	if errMsg.Valid {
		j.ErrorMessage = errMsg.String
	}
	return &j, nil
}

func (s *JobService) ListJobs(limit, offset int) ([]models.TranslationJob, int, error) {
	var total int
	s.DB.QueryRow("SELECT COUNT(*) FROM translation_jobs").Scan(&total)

	rows, err := s.DB.Query(`
		SELECT id, status, file_path, target_lang, preset, model, total_lines, processed_lines, prompt_tokens, completion_tokens, cost_usd, error_message, created_at, updated_at 
		FROM translation_jobs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []models.TranslationJob
	for rows.Next() {
		var j models.TranslationJob
		var errMsg sql.NullString
		if err := rows.Scan(&j.ID, &j.Status, &j.FilePath, &j.TargetLang, &j.Preset, &j.Model, &j.TotalLines, &j.ProcessedLines, &j.PromptTokens, &j.CompletionTokens, &j.CostUSD, &errMsg, &j.CreatedAt, &j.UpdatedAt); err == nil {
			if errMsg.Valid {
				j.ErrorMessage = errMsg.String
			}
			jobs = append(jobs, j)
		}
	}
	if jobs == nil {
		jobs = []models.TranslationJob{}
	}
	return jobs, total, nil
}
