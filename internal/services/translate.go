package services

import (
	"bakasub-backend/internal/ai"
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/parser"
	"bakasub-backend/internal/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type LLMProvider interface {
	TranslateText(text string, model string, apiKey string, sourceLangName string, targetLangName string, systemPrompt string) (string, int, int, error)
	GetModelPricing(modelID string) (float64, float64, error)
}

type TranslationFileSystemProvider interface {
	ReadFile(path string) (string, error)
	SaveFile(path string, content string) error
}

type TranslatorService struct {
	LLM LLMProvider
	FS  TranslationFileSystemProvider
	DB  *sql.DB
	Job *JobService
}

func NewTranslatorService(llm LLMProvider, fs TranslationFileSystemProvider, db *sql.DB, job *JobService) *TranslatorService {
	return &TranslatorService{
		LLM: llm,
		FS:  fs,
		DB:  db,
		Job: job,
	}
}

var separatorRegex = regexp.MustCompile(`(?:\\N|\s)*---NEXT---(?:\\N|\s)*`)
var looseSeparatorRegex = regexp.MustCompile(`(?:\\N|\n)\s*-{3,}\s*(?:\\N|\n)`)
var langSuffixRe = regexp.MustCompile(`[._-]([a-zA-Z]{2,3}(?:-[a-zA-Z]{2,3})?)$`)

func (s *TranslatorService) PreFlight(inputPath string, model string, targetLangCode string, presetAlias string, removeSDH bool, contextData string) (*models.JobEstimate, error) {
	rawText, err := s.FS.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file for preflight: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	var blocks []models.SubtitleBlock

	switch ext {
	case ".ass", ".ssa":
		_, blocks = parser.ParseASS(rawText)
	case ".vtt":
		_, blocks = parser.ParseVTT(rawText)
	default:
		blocks = parser.ParseToBlocks(rawText)
	}

	if removeSDH && ext != ".ass" {
		blocks = parser.RemoveSDH(blocks)
	}

	totalLines := len(blocks)
	if totalLines == 0 {
		return nil, fmt.Errorf("no valid subtitle found to estimate")
	}

	ctxConfig, err := s.getTranslationContext(targetLangCode, presetAlias, contextData)
	if err != nil {
		return nil, fmt.Errorf("error fetching translation context for preflight: %w", err)
	}

	sourceLanguage := s.detectSourceLanguage(inputPath)

	uncachedIndices, _ := s.applyTranslationMemory(blocks, targetLangCode, presetAlias)

	cachedLines := totalLines - len(uncachedIndices)
	linesToTranslate := len(uncachedIndices)

	charCount := 0
	for _, idx := range uncachedIndices {
		charCount += len(blocks[idx].Text)
	}

	estInputTokens := 0
	estOutputTokens := 0
	totalBatches := 0

	if linesToTranslate > 0 {
		batches := s.createBatches(blocks, uncachedIndices, ctxConfig.MaxChars)
		totalBatches = len(batches)

		estInputTokens = (charCount / 4) + (500 * totalBatches)
		estOutputTokens = charCount / 4
	}

	promptPrice, completionPrice, err := s.LLM.GetModelPricing(model)
	if err != nil {
		promptPrice = 0
		completionPrice = 0
	}

	estCost := (float64(estInputTokens) * promptPrice) + (float64(estOutputTokens) * completionPrice)

	var presetName string
	_ = s.DB.QueryRow("SELECT name FROM translation_presets WHERE alias = $1", presetAlias).Scan(&presetName)

	return &models.JobEstimate{
		TotalLines:       totalLines,
		CachedLines:      cachedLines,
		LinesToTranslate: linesToTranslate,
		TotalBatches:     totalBatches,
		EstimatedTokens:  estInputTokens + estOutputTokens,
		EstimatedCostUSD: estCost,
		IsFreeModel:      ai.IsModelFree(model),
		SystemPrompt:     ctxConfig.SystemPrompt,
		TargetLanguage:   ctxConfig.TargetLangName,
		SourceLanguage:   sourceLanguage,
		PresetName:       presetName,
		BatchSize:        ctxConfig.MaxChars,
	}, nil
}

func (s *TranslatorService) ProcessSubtitleFile(jobID, inputPath, model, outputPath, apiKey, targetLangCode, presetAlias string, removeSDH bool, contextData string) error {
	if apiKey == "" {
		return fmt.Errorf("apiKey is empty before starting translation")
	}

	utils.LogInfo("translate", "info", "Starting translation process", map[string]any{"job_id": jobID, "file": filepath.Base(inputPath)})
	utils.SendSSE("info", "translate", "Starting subtitle translation...", map[string]any{"job_id": jobID})

	ctxConfig, err := s.getTranslationContext(targetLangCode, presetAlias, contextData)
	if err != nil {
		return err
	}

	ctxConfig.SourceLangName = s.detectSourceLanguage(inputPath)

	promptPrice, completionPrice, err := s.LLM.GetModelPricing(model)
	if err == nil {
		ctxConfig.PromptPrice = promptPrice
		ctxConfig.CompletionPrice = completionPrice
	}

	rawText, err := s.FS.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	var blocks []models.SubtitleBlock
	var assDoc *parser.ASSDocument
	var vttHeader string

	switch ext {
	case ".ass", ".ssa":
		assDoc, blocks = parser.ParseASS(rawText)
	case ".vtt":
		vttHeader, blocks = parser.ParseVTT(rawText)
	default:
		blocks = parser.ParseToBlocks(rawText)
	}

	if removeSDH && ext != ".ass" {
		blocks = parser.RemoveSDH(blocks)
	}
	if len(blocks) == 0 {
		return fmt.Errorf("no valid subtitle found after parsing")
	}

	s.Job.UpdateTotalLines(jobID, len(blocks))

	uncachedIndices, originalTexts := s.applyTranslationMemory(blocks, targetLangCode, presetAlias)
	cachedCount := len(blocks) - len(uncachedIndices)

	if cachedCount > 0 {
		s.Job.IncrementProgress(jobID, cachedCount, 0, 0, 0)
		s.Job.SetCachedLines(jobID, cachedCount)
	}

	batches := s.createBatches(blocks, uncachedIndices, ctxConfig.MaxChars)

	if len(batches) > 0 {
		err = s.processTranslationBatches(jobID, batches, blocks, originalTexts, model, apiKey, ctxConfig)
		if err != nil {
			return err
		}
	}

	err = s.buildAndSaveSubtitle(ext, outputPath, blocks, assDoc, vttHeader, targetLangCode)
	if err != nil {
		return err
	}

	utils.LogInfo("translate", "success", "Translation finished", map[string]any{"job_id": jobID})

	successData := map[string]any{
		"job_id": jobID,
		"output": filepath.Base(outputPath),
	}
	if job, err := s.Job.GetJob(jobID); err == nil {
		successData["total_lines"] = job.TotalLines
		successData["processed_lines"] = job.ProcessedLines
		successData["cached_lines"] = job.CachedLines
		successData["prompt_tokens"] = job.PromptTokens
		successData["completion_tokens"] = job.CompletionTokens
		successData["cost_usd"] = job.CostUSD
		successData["model"] = job.Model
		successData["created_at"] = job.CreatedAt
		successData["updated_at"] = job.UpdatedAt
	}
	utils.SendSSE("success", "translate", "Translation finished successfully!", successData)

	return nil
}

type TranslationContext struct {
	SystemPrompt    string
	MaxChars        int
	TargetLangName  string
	SourceLangName  string
	MaxRetries      int
	Concurrency     int
	BaseRetryDelay  int
	PresetAlias     string
	TargetLangCode  string
	PromptPrice     float64
	CompletionPrice float64
}

func (s *TranslatorService) getTranslationContext(targetLangCode, presetAlias, contextData string) (TranslationContext, error) {
	var ctx TranslationContext
	ctx.TargetLangCode = targetLangCode
	ctx.PresetAlias = presetAlias

	err := s.DB.QueryRow("SELECT system_prompt, batch_size FROM translation_presets WHERE alias = $1", presetAlias).Scan(&ctx.SystemPrompt, &ctx.MaxChars)
	if err != nil {
		return ctx, fmt.Errorf("error fetching preset '%s': %w", presetAlias, err)
	}

	if strings.TrimSpace(contextData) != "" {
		ctx.SystemPrompt += fmt.Sprintf("\n\nMEDIA CONTEXT, use this informations to understand the media content, characters names and locations, history, and other relevant details:\n%s", strings.TrimSpace(contextData))
	}

	err = s.DB.QueryRow("SELECT name FROM languages WHERE code = $1", targetLangCode).Scan(&ctx.TargetLangName)
	if err != nil {
		return ctx, fmt.Errorf("error fetching language '%s': %w", targetLangCode, err)
	}

	err = s.DB.QueryRow("SELECT concurrent_translations, max_retries, base_retry_delay FROM user_config LIMIT 1").
		Scan(&ctx.Concurrency, &ctx.MaxRetries, &ctx.BaseRetryDelay)
	if err != nil {
		ctx.Concurrency = 5
		ctx.MaxRetries = 3
		ctx.BaseRetryDelay = 2
	}

	return ctx, nil
}

func (s *TranslatorService) detectSourceLanguage(inputPath string) string {
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	matches := langSuffixRe.FindStringSubmatch(name)
	if len(matches) < 2 {
		return ""
	}

	alias := strings.ToLower(matches[1])
	var langName string

	query := `
		SELECT l.name 
		FROM language_mappings m
		JOIN languages l ON m.language_code = l.code
		WHERE m.alias = $1
	`
	err := s.DB.QueryRow(query, alias).Scan(&langName)
	if err != nil {
		return ""
	}

	return langName
}

func (s *TranslatorService) applyTranslationMemory(blocks []models.SubtitleBlock, targetLangCode, presetAlias string) ([]int, map[int]string) {
	originalTexts := make(map[int]string)
	if len(blocks) == 0 {
		return nil, originalTexts
	}

	// Compute all hashes upfront
	blockHashes := make([]string, len(blocks))
	for i, block := range blocks {
		originalTexts[i] = block.Text
		hashInput := block.Text + "|" + targetLangCode + "|" + presetAlias
		hashBytes := sha256.Sum256([]byte(hashInput))
		blockHashes[i] = hex.EncodeToString(hashBytes[:])
	}

	// Single batch query instead of N individual queries
	placeholders := make([]string, len(blockHashes))
	args := make([]any, len(blockHashes))
	for i, h := range blockHashes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = h
	}
	query := fmt.Sprintf("SELECT hash, translated_text FROM translation_memory WHERE hash IN (%s)", strings.Join(placeholders, ","))

	cachedMap := make(map[string]string)
	rows, err := s.DB.Query(query, args...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var h, t string
			if rows.Scan(&h, &t) == nil && t != "" {
				cachedMap[h] = t
			}
		}
	}

	var uncachedIndices []int
	for i := range blocks {
		if trans, ok := cachedMap[blockHashes[i]]; ok {
			blocks[i].Text = trans
		} else {
			uncachedIndices = append(uncachedIndices, i)
		}
	}

	return uncachedIndices, originalTexts
}

func (s *TranslatorService) createBatches(blocks []models.SubtitleBlock, uncachedIndices []int, maxChars int) [][]int {
	var batches [][]int
	var currentBatch []int
	currentChars := 0
	const separatorLen = 12

	for _, idx := range uncachedIndices {
		blockLen := len(blocks[idx].Text)

		if currentChars+blockLen > maxChars && len(currentBatch) > 0 {
			batches = append(batches, currentBatch)
			currentBatch = nil
			currentChars = 0
		}

		currentBatch = append(currentBatch, idx)
		currentChars += blockLen + separatorLen
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}

func (s *TranslatorService) processTranslationBatches(jobID string, batches [][]int, blocks []models.SubtitleBlock, originalTexts map[int]string, model, apiKey string, ctxConfig TranslationContext) error {
	isFree := ai.IsModelFree(model)

	if isFree {
		if ctxConfig.Concurrency > 1 {
			ctxConfig.Concurrency = 1
		}
		if ctxConfig.BaseRetryDelay < 10 {
			ctxConfig.BaseRetryDelay = 10
		}
		if ctxConfig.MaxRetries > 2 {
			ctxConfig.MaxRetries = 2
		}
		utils.LogInfo("translate", "warn", "Free model detected — using reduced concurrency (1), extended backoff (10s) and lower max retries (2)", map[string]any{
			"job_id": jobID,
			"model":  model,
		})
		utils.SendSSE("warn", "translate", "Free model detected — using reduced concurrency (1) and extended backoff (10s). Free models have strict rate limits; consider using a paid model for large files.", map[string]any{
			"job_id": jobID,
			"model":  model,
		})
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var translationErr error

	limit := make(chan struct{}, ctxConfig.Concurrency)
	totalBatches := len(batches)
	completedBatches := 0

	const separator = "\n---NEXT---\n"

	for _, batchIndices := range batches {
		wg.Add(1)
		limit <- struct{}{}

		go func(indices []int) {
			defer wg.Done()
			defer func() { <-limit }()

			mu.Lock()
			if translationErr != nil {
				mu.Unlock()
				return
			}
			mu.Unlock()

			var texts []string
			for _, idx := range indices {
				texts = append(texts, blocks[idx].Text)
			}
			joinedText := strings.Join(texts, separator)

			var translatedText string
			var pTokens, cTokens int
			var err error
			backoff := time.Duration(ctxConfig.BaseRetryDelay) * time.Second

			for attempt := 1; attempt <= ctxConfig.MaxRetries; attempt++ {
				translatedText, pTokens, cTokens, err = s.LLM.TranslateText(joinedText, model, apiKey, ctxConfig.SourceLangName, ctxConfig.TargetLangName, ctxConfig.SystemPrompt)
				if err == nil {
					break
				}

				var apiErr *ai.APIError
				if errors.As(err, &apiErr) {
					if !apiErr.Retryable {
						utils.LogInfo("translate", "error", "Fatal API error, aborting without retry", map[string]any{
							"status": apiErr.StatusCode,
							"error":  apiErr.Message,
						})
						break
					}

					if apiErr.StatusCode == http.StatusTooManyRequests && isFree {
						err = fmt.Errorf("[Free Model] rate limit exceeded (HTTP 429). Free models have very strict rate limits — try reducing batch size, using a paid model, or waiting a few minutes. Details: %s", apiErr.Message)
						utils.LogInfo("translate", "error", "Free model rate limited, failing immediately", map[string]any{
							"model": model,
							"error": apiErr.Message,
						})
						utils.SendSSE("warn", "translate", "Free model rate limit hit (429). Free models have strict limits — consider using a paid model.", map[string]any{
							"job_id": jobID,
							"model":  model,
						})
						break
					}

					if apiErr.RetryAfter > 0 {
						backoff = apiErr.RetryAfter
					}
				}

				utils.LogInfo("translate", "warn", "Batch translation failed, retrying...", map[string]any{"attempt": attempt, "error": err.Error()})
				if attempt < ctxConfig.MaxRetries {
					time.Sleep(backoff)
					backoff *= 2
				}
			}

			if err != nil {
				mu.Lock()
				if translationErr == nil {
					translationErr = fmt.Errorf("error translating batch starting at line %d after %d attempts: %w", indices[0], ctxConfig.MaxRetries, err)
				}
				mu.Unlock()
				return
			}

			translatedLines := separatorRegex.Split(translatedText, -1)
			for len(translatedLines) > 0 && strings.TrimSpace(translatedLines[0]) == "" {
				translatedLines = translatedLines[1:]
			}
			for len(translatedLines) > 0 && strings.TrimSpace(translatedLines[len(translatedLines)-1]) == "" {
				translatedLines = translatedLines[:len(translatedLines)-1]
			}

			if len(translatedLines) != len(indices) {
				utils.LogInfo("translate", "warn", "Separator count mismatch after strict split", map[string]any{
					"expected": len(indices),
					"got":      len(translatedLines),
				})

				// Try loose separator: catches LLM deviations like \N---\N
				altLines := looseSeparatorRegex.Split(translatedText, -1)
				for len(altLines) > 0 && strings.TrimSpace(altLines[0]) == "" {
					altLines = altLines[1:]
				}
				for len(altLines) > 0 && strings.TrimSpace(altLines[len(altLines)-1]) == "" {
					altLines = altLines[:len(altLines)-1]
				}
				if len(altLines) == len(indices) {
					translatedLines = altLines
				} else {
					// Final fallback: split by blank lines without regex replacement
					fallbackLines := strings.Split(strings.TrimSpace(translatedText), "\n\n")
					translatedLines = make([]string, len(indices))
					for i := range fallbackLines {
						if i < len(indices) {
							translatedLines[i] = fallbackLines[i]
						}
					}
					utils.LogInfo("translate", "warn", "Used fallback blank-line splitting", map[string]any{
						"expected":     len(indices),
						"got_loose":    len(altLines),
						"got_fallback": len(fallbackLines),
					})
				}
			}

			mu.Lock()
			for i, idx := range indices {
				if i < len(translatedLines) {
					finalTranslatedText := strings.TrimSpace(translatedLines[i])
					blocks[idx].Text = finalTranslatedText

					origText := originalTexts[idx]
					hashInput := origText + "|" + ctxConfig.TargetLangCode + "|" + ctxConfig.PresetAlias
					hashBytes := sha256.Sum256([]byte(hashInput))
					hashStr := hex.EncodeToString(hashBytes[:])

					s.DB.Exec(`
                        INSERT INTO translation_memory (hash, source_text, translated_text, target_lang, preset) 
                        VALUES ($1, $2, $3, $4, $5) 
                        ON CONFLICT (hash) DO NOTHING`,
						hashStr, origText, finalTranslatedText, ctxConfig.TargetLangCode, ctxConfig.PresetAlias,
					)
				}
			}

			batchCost := (float64(pTokens) * ctxConfig.PromptPrice) + (float64(cTokens) * ctxConfig.CompletionPrice)
			s.Job.IncrementProgress(jobID, len(indices), pTokens, cTokens, batchCost)

			completedBatches++
			percent := (completedBatches * 100) / totalBatches
			utils.SendSSE("progress", "translate", fmt.Sprintf("Translating batch %d of %d (%d%%)", completedBatches, totalBatches, percent), map[string]any{
				"job_id":  jobID,
				"current": completedBatches,
				"total":   totalBatches,
				"percent": percent,
			})
			mu.Unlock()

		}(batchIndices)
	}

	wg.Wait()
	return translationErr
}

func (s *TranslatorService) buildAndSaveSubtitle(ext, outputPath string, blocks []models.SubtitleBlock, assDoc *parser.ASSDocument, vttHeader string, targetLang string) error {
	var outputText string
	switch ext {
	case ".ass", ".ssa":
		newTitle := fmt.Sprintf("Title: [BakaSub-AI] %s", targetLang)
		if strings.Contains(assDoc.Header, "[Script Info]") {
			assDoc.Header = strings.Replace(assDoc.Header, "[Script Info]", "[Script Info]\n"+newTitle, 1)
		} else {
			assDoc.Header = "[Script Info]\n" + newTitle + "\n" + assDoc.Header
		}
		outputText = parser.BuildASS(assDoc, blocks)
	case ".vtt":
		outputText = parser.BuildVTT(vttHeader, blocks)
	default:
		outputText = parser.BuildString(blocks)
	}

	return s.FS.SaveFile(outputPath, outputText)
}
