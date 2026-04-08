package parser

import (
	"bakasub-backend/internal/models"
	"fmt"
	"strings"
)

func ParseVTT(rawText string) (string, []models.SubtitleBlock) {
	normalized := strings.ReplaceAll(rawText, "\r\n", "\n")
	rawBlocks := strings.Split(normalized, "\n\n")

	var headerBuilder strings.Builder
	var parsedBlocks []models.SubtitleBlock

	for i, block := range rawBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		if i == 0 && strings.HasPrefix(block, "WEBVTT") {
			headerBuilder.WriteString(block + "\n\n")
			continue
		}

		lines := strings.Split(block, "\n")
		if len(lines) < 2 {
			continue
		}

		var id, time string
		var textLines []string

		if strings.Contains(lines[0], "-->") {
			time = lines[0]
			textLines = lines[1:]
		} else if len(lines) >= 3 && strings.Contains(lines[1], "-->") {
			id = lines[0]
			time = lines[1]
			textLines = lines[2:]
		} else {
			continue
		}

		var cleanTexts []string
		var formats []models.LineFormat

		for _, line := range textLines {
			prefix := prefixRe.FindString(line)
			remainder := line[len(prefix):]
			suffix := suffixRe.FindString(remainder)
			cleanText := remainder[:len(remainder)-len(suffix)]

			if cleanText != "" {
				cleanTexts = append(cleanTexts, cleanText)
				formats = append(formats, models.LineFormat{
					Prefix: prefix,
					Suffix: suffix,
				})
			}
		}

		if len(cleanTexts) == 0 {
			continue
		}

		parsedBlocks = append(parsedBlocks, models.SubtitleBlock{
			ID:         id,
			Time:       time,
			Text:       strings.Join(cleanTexts, "\n"),
			Formatting: formats,
		})
	}
	return headerBuilder.String(), parsedBlocks
}

func BuildVTT(header string, blocks []models.SubtitleBlock) string {
	var builder strings.Builder

	if header == "" {
		builder.WriteString("WEBVTT\n\n")
	} else {
		builder.WriteString(header)
	}

	for _, block := range blocks {
		if block.ID != "" {
			builder.WriteString(fmt.Sprintf("%s\n", block.ID))
		}
		builder.WriteString(fmt.Sprintf("%s\n", block.Time))

		translatedLines := strings.Split(block.Text, "\n")
		for i, textLine := range translatedLines {
			prefix := ""
			suffix := ""

			if i < len(block.Formatting) {
				prefix = block.Formatting[i].Prefix
				suffix = block.Formatting[i].Suffix
			}
			if i == len(translatedLines)-1 && i >= len(block.Formatting) && len(block.Formatting) > 0 {
				suffix = block.Formatting[len(block.Formatting)-1].Suffix
			}

			builder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, strings.TrimSpace(textLine), suffix))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
