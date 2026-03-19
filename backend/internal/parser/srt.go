package parser

import (
	"bakasub-backend/internal/models"
	"fmt"
	"regexp"
	"strings"
)

var prefixRe = regexp.MustCompile(`^(?:<[^>]+>|\{\\[^}]+\}|\s)+`)
var suffixRe = regexp.MustCompile(`(?:<[^>]+>|\{\\[^}]+\}|\s)+$`)

func ParseToBlocks(rawText string) []models.SubtitleBlock {
	normalizedSrt := strings.ReplaceAll(rawText, "\r\n", "\n")
	rawBlocks := strings.Split(normalizedSrt, "\n\n")

	var parsedBlocks []models.SubtitleBlock

	for _, block := range rawBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		lines := strings.Split(block, "\n")
		if len(lines) < 3 {
			continue
		}

		id := strings.TrimSpace(lines[0])
		time := strings.TrimSpace(lines[1])
		textLines := lines[2:]

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
	return parsedBlocks
}

func BuildString(blocks []models.SubtitleBlock) string {
	var builder strings.Builder

	for _, block := range blocks {
		builder.WriteString(fmt.Sprintf("%s\n%s\n", block.ID, block.Time))

		translatedLines := strings.Split(block.Text, "\n")

		for i, textLine := range translatedLines {
			prefix := ""
			suffix := ""

			if i < len(block.Formatting) {
				prefix = block.Formatting[i].Prefix
				suffix = block.Formatting[i].Suffix
			}

			if i == len(translatedLines)-1 && len(block.Formatting) > 0 {
				suffix = block.Formatting[len(block.Formatting)-1].Suffix
			}

			textLine = strings.TrimSpace(textLine)

			builder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, textLine, suffix))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
