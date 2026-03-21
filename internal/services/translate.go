package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/parser"
	"bakasub-backend/internal/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type LLMProvider interface {
	TranslateText(text string, model string, apiKey string, targetLangName string, systemPrompt string) (string, error)
}

type TranslationFileSystemProvider interface {
	ReadFile(path string) (string, error)
	SaveFile(path string, content string) error
}

type TranslatorService struct {
	LLM LLMProvider
	FS  TranslationFileSystemProvider
	DB  *sql.DB
}

func NewTranslatorService(llm LLMProvider, fs TranslationFileSystemProvider, db *sql.DB) *TranslatorService {
	return &TranslatorService{
		LLM: llm,
		FS:  fs,
		DB:  db,
	}
}

var separatorRegex = regexp.MustCompile(`\s*---NEXT---\s*`)

func (s *TranslatorService) ProcessSubtitleFile(inputPath, model, outputPath, apiKey, targetLangCode, presetAlias string, removeSDH bool, contextData string) error {
	if apiKey == "" {
		return fmt.Errorf("apiKey is empty before starting translation")
	}

	utils.LogInfo("translate", "info", "Starting translation process", map[string]any{
		"file": filepath.Base(inputPath),
		"lang": targetLangCode,
	})
	utils.SendSSE("info", "translate", "Starting subtitle translation...", map[string]any{"file": filepath.Base(inputPath)})

	ctxConfig, err := s.getTranslationContext(targetLangCode, presetAlias, contextData)
	if err != nil {
		utils.LogError("translate", "Failed to fetch translation context", map[string]any{"error": err.Error()})
		utils.SendSSE("error", "translate", "Failed to load translation settings.", nil)
		return err
	}

	rawText, err := s.FS.ReadFile(inputPath)
	if err != nil {
		utils.LogError("translate", "Failed to read subtitle file", map[string]any{"error": err.Error()})
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
		return fmt.Errorf("no valid subtitle found after parsing and filtering")
	}

	uncachedIndices, originalTexts := s.applyTranslationMemory(blocks, targetLangCode, presetAlias)

	utils.LogInfo("translate", "info", "Memory cache applied", map[string]any{
		"total_blocks":  len(blocks),
		"cached_blocks": len(blocks) - len(uncachedIndices),
		"to_translate":  len(uncachedIndices),
	})

	batches := s.createBatches(blocks, uncachedIndices, ctxConfig.MaxChars)

	if len(batches) > 0 {
		utils.SendSSE("info", "translate", "Communicating with AI model...", map[string]any{"total_batches": len(batches)})

		err = s.processTranslationBatches(batches, blocks, originalTexts, model, apiKey, ctxConfig)
		if err != nil {
			utils.LogError("translate", "Batch processing failed", map[string]any{"error": err.Error()})
			utils.SendSSE("error", "translate", "Translation process failed.", nil)
			return err
		}
	}

	err = s.buildAndSaveSubtitle(ext, outputPath, blocks, assDoc, vttHeader)
	if err != nil {
		utils.LogError("translate", "Failed to save final subtitle", map[string]any{"error": err.Error()})
		utils.SendSSE("error", "translate", "Failed to save output file.", nil)
		return err
	}

	utils.LogInfo("translate", "success", "Translation finished successfully", map[string]any{"output": outputPath})
	utils.SendSSE("success", "translate", "Translation finished successfully!", map[string]any{"output": filepath.Base(outputPath)})

	return nil
}

type TranslationContext struct {
	SystemPrompt   string
	MaxChars       int
	TargetLangName string
	MaxRetries     int
	Concurrency    int
	BaseRetryDelay int
	PresetAlias    string
	TargetLangCode string
}

func (s *TranslatorService) getTranslationContext(targetLangCode, presetAlias, contextData string) (TranslationContext, error) {
	var ctx TranslationContext
	ctx.TargetLangCode = targetLangCode
	ctx.PresetAlias = presetAlias

	err := s.DB.QueryRow("SELECT system_prompt, batch_size FROM translation_presets WHERE alias = ?", presetAlias).Scan(&ctx.SystemPrompt, &ctx.MaxChars)
	if err != nil {
		return ctx, fmt.Errorf("error fetching preset '%s': %w", presetAlias, err)
	}

	if strings.TrimSpace(contextData) != "" {
		ctx.SystemPrompt += fmt.Sprintf("\n\nMEDIA CONTEXT, use this informations to understand the media content, characters names and locations, history, and other relevant details:\n%s", strings.TrimSpace(contextData))
	}

	err = s.DB.QueryRow("SELECT name FROM languages WHERE code = ?", targetLangCode).Scan(&ctx.TargetLangName)
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

func (s *TranslatorService) applyTranslationMemory(blocks []models.SubtitleBlock, targetLangCode, presetAlias string) ([]int, map[int]string) {
	var uncachedIndices []int
	originalTexts := make(map[int]string)

	for i, block := range blocks {
		originalTexts[i] = block.Text

		hashInput := block.Text + "|" + targetLangCode + "|" + presetAlias
		hashBytes := sha256.Sum256([]byte(hashInput))
		hashStr := hex.EncodeToString(hashBytes[:])

		var cachedTranslation string
		err := s.DB.QueryRow("SELECT translated_text FROM translation_memory WHERE hash = ?", hashStr).Scan(&cachedTranslation)

		if err == nil && cachedTranslation != "" {
			blocks[i].Text = cachedTranslation // Cache Hit
		} else {
			uncachedIndices = append(uncachedIndices, i) // Cache Miss
		}
	}

	return uncachedIndices, originalTexts
}

func (s *TranslatorService) createBatches(blocks []models.SubtitleBlock, uncachedIndices []int, maxChars int) [][]int {
	var batches [][]int
	var currentBatch []int
	currentChars := 0
	const separatorLen = 12 // "\n---NEXT---\n" len

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

func (s *TranslatorService) processTranslationBatches(batches [][]int, blocks []models.SubtitleBlock, originalTexts map[int]string, model, apiKey string, ctxConfig TranslationContext) error {
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
			var err error
			backoff := time.Duration(ctxConfig.BaseRetryDelay) * time.Second

			for attempt := 1; attempt <= ctxConfig.MaxRetries; attempt++ {
				translatedText, err = s.LLM.TranslateText(joinedText, model, apiKey, ctxConfig.TargetLangName, ctxConfig.SystemPrompt)
				if err == nil {
					break // Success
				}

				utils.LogInfo("translate", "warn", "Batch translation failed, retrying...", map[string]any{
					"attempt": attempt,
					"error":   err.Error(),
				})

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
			if len(translatedLines) > 0 && strings.TrimSpace(translatedLines[len(translatedLines)-1]) == "" {
				translatedLines = translatedLines[:len(translatedLines)-1]
			}

			if len(translatedLines) != len(indices) {
				cleanText := separatorRegex.ReplaceAllString(translatedText, "\n\n")
				fallbackLines := strings.Split(strings.TrimSpace(cleanText), "\n\n")
				translatedLines = make([]string, len(indices))
				for i := range fallbackLines {
					if i < len(indices) {
						translatedLines[i] = fallbackLines[i]
					}
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

					_, _ = s.DB.Exec(`
						INSERT OR IGNORE INTO translation_memory (hash, source_text, translated_text, target_lang, preset) 
						VALUES (?, ?, ?, ?, ?)`,
						hashStr, origText, finalTranslatedText, ctxConfig.TargetLangCode, ctxConfig.PresetAlias,
					)
				}
			}

			completedBatches++
			percent := (completedBatches * 100) / totalBatches
			utils.SendSSE("progress", "translate", fmt.Sprintf("Translating batch %d of %d (%d%%)", completedBatches, totalBatches, percent), map[string]any{
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

func (s *TranslatorService) buildAndSaveSubtitle(ext, outputPath string, blocks []models.SubtitleBlock, assDoc *parser.ASSDocument, vttHeader string) error {
	var outputText string
	switch ext {
	case ".ass", ".ssa":
		outputText = parser.BuildASS(assDoc, blocks)
	case ".vtt":
		outputText = parser.BuildVTT(vttHeader, blocks)
	default:
		outputText = parser.BuildString(blocks)
	}

	return s.FS.SaveFile(outputPath, outputText)
}
