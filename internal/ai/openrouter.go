package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

type ModelInfo struct {
	ID      string       `json:"id"`
	Pricing ModelPricing `json:"pricing"`
}

type ModelsResponse struct {
	Data []ModelInfo `json:"data"`
}

type pricingCacheEntry struct {
	Prompt     float64
	Completion float64
	Expires    time.Time
}

type OpenRouterService struct {
	pricingCache map[string]pricingCacheEntry
	mu           sync.RWMutex
}

func NewOpenRouterService() *OpenRouterService {
	return &OpenRouterService{
		pricingCache: make(map[string]pricingCacheEntry),
	}
}

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

func (s *OpenRouterService) GetModelPricing(modelID string) (float64, float64, error) {
	s.mu.RLock()
	entry, exists := s.pricingCache[modelID]
	s.mu.RUnlock()

	if exists && time.Now().Before(entry.Expires) {
		return entry.Prompt, entry.Completion, nil
	}

	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("failed to fetch models, status code: %d", resp.StatusCode)
	}

	var res ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, 0, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var promptCost, completionCost float64
	found := false

	for _, m := range res.Data {
		p, _ := strconv.ParseFloat(m.Pricing.Prompt, 64)
		c, _ := strconv.ParseFloat(m.Pricing.Completion, 64)

		s.pricingCache[m.ID] = pricingCacheEntry{
			Prompt:     p,
			Completion: c,
			Expires:    time.Now().Add(1 * time.Hour),
		}

		if m.ID == modelID {
			promptCost = p
			completionCost = c
			found = true
		}
	}

	if !found {
		return 0, 0, fmt.Errorf("model %s not found in OpenRouter", modelID)
	}

	return promptCost, completionCost, nil
}

func (s *OpenRouterService) TranslateText(text string, model string, apiKey string, sourceLangName string, targetLangName string, systemPrompt string) (string, int, int, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	taskInstruction := fmt.Sprintf("Translate the input to %s.", targetLangName)
	if sourceLangName != "" {
		taskInstruction = fmt.Sprintf("Translate the input from %s to %s.", sourceLangName, targetLangName)
	}
	if strings.Contains(text, "---NEXT---") {
		taskInstruction += " The input contains multiple subtitle items separated by ---NEXT---. Translate each item individually and preserve the ---NEXT--- separator between each translated item exactly as it appears in the input. Do not merge, reorder, or skip any items."
	}

	reqBody := RequestBody{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("%s\n\nTask: %s", systemPrompt, taskInstruction),
			},
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("OpenRouter API error (Status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var resBody ResponseBody
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		return "", 0, 0, fmt.Errorf("error unmarshaling response: %v | Body: %s", err, string(bodyBytes))
	}

	if len(resBody.Choices) == 0 || resBody.Choices[0].Message.Content == "" {
		fmt.Printf("--- WARNING: Empty AI Response ---\nFull Body: %s\n-------------------\n", string(bodyBytes))
		return "", 0, 0, fmt.Errorf("no translation found (possible content filter)")
	}

	return resBody.Choices[0].Message.Content, resBody.Usage.PromptTokens, resBody.Usage.CompletionTokens, nil
}
