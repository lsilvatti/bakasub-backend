package parser

import (
	"bakasub-backend/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var karaokeRe = regexp.MustCompile(`(?i)\{\\[k][fo]?[0-9.]+[^}]*\}`)
var prefixReASS = regexp.MustCompile(`^(?:\{[^}]+\}|\s)+`)
var suffixReASS = regexp.MustCompile(`(?:\{[^}]+\}|\s)+$`)
var nonTranslatableStyleRe = regexp.MustCompile(`(?i)\b(?:romaji|kanji|song|sign|op|ed)\b`)

type ASSDocument struct {
	Header string
	Lines  []*ASSLine
}

type ASSLine struct {
	IsTranslatable bool
	Raw            string
	Meta           string
	Prefix         string
	Suffix         string
	CleanText      string
	Hash           string
}

func ParseASS(rawText string) (*ASSDocument, []models.SubtitleBlock) {
	doc := &ASSDocument{}
	var headerBuilder strings.Builder
	var blocks []models.SubtitleBlock

	uniqueBlocks := make(map[string]bool)

	lines := strings.Split(strings.ReplaceAll(rawText, "\r\n", "\n"), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Comment:") {
			doc.Lines = append(doc.Lines, &ASSLine{IsTranslatable: false, Raw: line})
			continue
		}

		if !strings.HasPrefix(line, "Dialogue:") {
			if strings.HasPrefix(line, "Title:") {
				continue
			}
			headerBuilder.WriteString(line + "\n")
			continue
		}

		parts := strings.SplitN(line, ",", 10)
		if len(parts) < 10 {
			headerBuilder.WriteString(line + "\n")
			continue
		}

		meta := strings.Join(parts[:9], ",") + ","
		text := parts[9]
		style := strings.ToLower(parts[3])

		isNonTranslatable := nonTranslatableStyleRe.MatchString(style)
		hasKaraoke := karaokeRe.MatchString(text)

		if isNonTranslatable || hasKaraoke {
			doc.Lines = append(doc.Lines, &ASSLine{
				IsTranslatable: false,
				Raw:            line,
			})
			continue
		}

		textLines := strings.Split(strings.ReplaceAll(text, "\\n", "\\N"), "\\N")
		var cleanTexts []string
		var currentPrefix, currentSuffix string

		for i, tl := range textLines {
			p := prefixReASS.FindString(tl)
			rem := tl[len(p):]
			s := suffixReASS.FindString(rem)
			c := rem[:len(rem)-len(s)]

			if i == 0 {
				currentPrefix = p
			}
			if i == len(textLines)-1 {
				currentSuffix = s
			}
			if c != "" {
				cleanTexts = append(cleanTexts, c)
			}
		}

		cleanJoined := strings.Join(cleanTexts, "\n")
		if strings.TrimSpace(cleanJoined) == "" {
			doc.Lines = append(doc.Lines, &ASSLine{IsTranslatable: false, Raw: line})
			continue
		}

		hashInput := parts[1] + parts[2] + cleanJoined
		hashBytes := sha256.Sum256([]byte(hashInput))
		hashStr := hex.EncodeToString(hashBytes[:])

		doc.Lines = append(doc.Lines, &ASSLine{
			IsTranslatable: true,
			Meta:           meta,
			Prefix:         currentPrefix,
			Suffix:         currentSuffix,
			CleanText:      cleanJoined,
			Hash:           hashStr,
		})

		if !uniqueBlocks[hashStr] {
			uniqueBlocks[hashStr] = true
			blocks = append(blocks, models.SubtitleBlock{
				ID:   hashStr,
				Time: meta,
				Text: cleanJoined,
			})
		}
	}

	doc.Header = headerBuilder.String()
	return doc, blocks
}

func BuildASS(doc *ASSDocument, translatedBlocks []models.SubtitleBlock) string {
	translations := make(map[string]string)
	for _, b := range translatedBlocks {
		translations[b.ID] = b.Text
	}

	var builder strings.Builder
	builder.WriteString(doc.Header)

	for _, line := range doc.Lines {
		if !line.IsTranslatable {
			builder.WriteString(line.Raw + "\n")
			continue
		}

		translatedText, exists := translations[line.Hash]
		if !exists {
			translatedText = line.CleanText
		}

		translatedText = strings.ReplaceAll(translatedText, "\n", "\\N")

		finalLine := fmt.Sprintf("%s%s%s%s\n", line.Meta, line.Prefix, translatedText, line.Suffix)
		builder.WriteString(finalLine)
	}

	return builder.String()
}
