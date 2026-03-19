package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func TranslateText(text string, model string, apiKey string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"
	reqBody := RequestBody{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "Você é um tradutor profissional de legendas. Traduza do Inglês para o Português do Brasil. Responda EXCLUSIVAMENTE com o texto traduzido. Não adicione aspas ou formatação.",
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resBody ResponseBody
	if err := json.Unmarshal(bodyBytes, &resBody); err != nil {
		return "", err
	}

	if len(resBody.Choices) > 0 {
		return resBody.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("nenhuma tradução encontrada")
}
