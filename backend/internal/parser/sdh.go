package parser

import (
	"bakasub-backend/internal/models"
	"regexp"
	"strings"
)

var sdhRe = regexp.MustCompile((`(?i)\[.*?\]|\(.*?\)|♪+`))

func RemoveSDH(blocks []models.SubtitleBlock) []models.SubtitleBlock {
	var cleanedBlocks []models.SubtitleBlock

	for _, block := range blocks {
		lines := strings.Split(block.Text, "\n")
		var newLines []string
		var newFormatting []models.LineFormat

		for i, line := range lines {
			cleanedLine := sdhRe.ReplaceAllString(line, "")
			cleanedLine = strings.TrimSpace(cleanedLine)

			if cleanedLine != "" {
				newLines = append(newLines, cleanedLine)

				if i < len(block.Formatting) {
					newFormatting = append(newFormatting, block.Formatting[i])
				}
			}
		}

		if len(newLines) == 0 {
			continue
		}

		block.Text = strings.Join(newLines, "\n")
		block.Formatting = newFormatting
		cleanedBlocks = append(cleanedBlocks, block)
	}

	return cleanedBlocks
}
