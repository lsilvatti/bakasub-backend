package parser

import (
	"bakasub-backend/internal/models"
	"fmt"
	"strings"
)

func ParseToBlocks(rawText string) []models.SubtitleBlock {
	normalizedSrt := strings.ReplaceAll(rawText, "\r\n", "\n")
	rawBlocks := strings.Split(normalizedSrt, "\n\n")

	var parsedBlocks []models.SubtitleBlock

	for _, block := range rawBlocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		lines := strings.Split(block, "\n")
		if len(lines) < 3 {
			continue
		}
		parsedBlocks = append(parsedBlocks, models.SubtitleBlock{
			ID:   lines[0],
			Time: lines[1],
			Text: strings.Join(lines[2:], "\n"),
		})
	}
	return parsedBlocks
}

func BuildString(blocks []models.SubtitleBlock) string {
	var builder strings.Builder
	for _, block := range blocks {
		builder.WriteString(fmt.Sprintf("%s\n%s\n%s\n\n", block.ID, block.Time, block.Text))
	}
	return builder.String()
}
