package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

type OpenRouterService struct{}

func NewOpenRouterService() *OpenRouterService {
	return &OpenRouterService{}
}

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

func (s *OpenRouterService) TranslateText(text string, model string, apiKey string, targetLangName string, systemPrompt string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	reqBody := RequestBody{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("%s Translate input to target language: %s.", systemPrompt, targetLangName),
			},
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenRouter API error (Status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var resBody ResponseBody
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v | Body: %s", err, string(bodyBytes))
	}

	if len(resBody.Choices) == 0 || resBody.Choices[0].Message.Content == "" {
		fmt.Printf("--- ALERTA: Resposta Vazia da IA ---\nBody Completo: %s\n-------------------\n", string(bodyBytes))
		return "", fmt.Errorf("no translation found (possible content filter)")
	}

	return resBody.Choices[0].Message.Content, nil
}
