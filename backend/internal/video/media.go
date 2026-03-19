package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// --- STRUCTS PARA PARSE DO MKVMERGE ---

type mkvMergeProperties struct {
	TrackName    string `json:"track_name"`
	Language     string `json:"language"`
	DefaultTrack bool   `json:"default_track"`
	ForcedTrack  bool   `json:"forced_track"`
}

type mkvMergeTrack struct {
	ID         int                `json:"id"`
	Type       string             `json:"type"`
	Codec      string             `json:"codec"`
	Properties mkvMergeProperties `json:"properties"`
}

type mkvMergeOutput struct {
	Tracks []mkvMergeTrack `json:"tracks"`
}

type SubtitleTrack struct {
	ID        int    `json:"id"`
	Language  string `json:"language"`
	Name      string `json:"name"`
	Codec     string `json:"codec"`
	IsDefault bool   `json:"isDefault"`
	IsForced  bool   `json:"isForced"`
}

func ScanSubtitles(videoPath string) ([]SubtitleTrack, error) {
	probeCmd := exec.Command("mkvmerge", "-J", videoPath)
	var out bytes.Buffer
	probeCmd.Stdout = &out

	if err := probeCmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao executar probe no arquivo: %w", err)
	}

	var info mkvMergeOutput
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON do mkvmerge: %w", err)
	}

	var subtitleTracks []SubtitleTrack

	for _, track := range info.Tracks {
		if track.Type == "subtitles" {
			trackName := track.Properties.TrackName
			if trackName == "" {
				trackName = fmt.Sprintf("Track %d", track.ID)
			}

			lang := track.Properties.Language
			if lang == "" {
				lang = "und"
			}

			cleanTrack := SubtitleTrack{
				ID:        track.ID,
				Language:  lang,
				Name:      trackName,
				Codec:     track.Codec,
				IsDefault: track.Properties.DefaultTrack,
				IsForced:  track.Properties.ForcedTrack,
			}

			subtitleTracks = append(subtitleTracks, cleanTrack)
		}
	}

	return subtitleTracks, nil
}

func ExtractSubtitle(videoPath string, subtitleId int) (string, error) {
	tracks, err := ScanSubtitles(videoPath)
	if err != nil || len(tracks) == 0 {
		return "", fmt.Errorf("falha ao ler as trilhas do vídeo ou nenhuma trilha encontrada: %v", err)
	}

	lang := "und"
	trackFound := false
	for _, t := range tracks {
		if t.ID == subtitleId {
			lang = t.Language
			trackFound = true
			break
		}
	}

	if !trackFound {
		return "", fmt.Errorf("trilha com ID %d não foi encontrada", subtitleId)
	}

	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	srtFilename := fmt.Sprintf("%s_%d_%s.srt", nameWithoutExt, subtitleId, lang)
	srtPath := filepath.Join(dir, srtFilename)

	fileExt := strings.ToLower(ext)

	if fileExt == ".mkv" {
		err = ExtractWithMKVToolnix(videoPath, subtitleId, srtPath)
	} else {
		err = ExtractWithFFmpeg(videoPath, subtitleId, srtPath)
	}

	if err != nil {
		return "", err
	}

	return srtPath, nil
}

func ExtractWithMKVToolnix(videoPath string, subtitleId int, srtPath string) error {
	extractCmd := exec.Command("mkvextract", videoPath, "tracks", fmt.Sprintf("%d:%s", subtitleId, srtPath))
	output, err := extractCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro rodando mkvextract: %w | log: %s", err, string(output))
	}
	return nil
}

func ExtractWithFFmpeg(videoPath string, subtitleId int, srtPath string) error {
	extractCmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-map", fmt.Sprintf("0:%d", subtitleId), srtPath)
	output, err := extractCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro extraindo com FFmpeg: %w | log: %s", err, string(output))
	}
	return nil
}

func MergeSubtitle(videoPath string, srtPath string, langCode string) (string, error) {
	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	outFilename := fmt.Sprintf("%s_bakasub_%s.mkv", nameWithoutExt, langCode)
	outVideoPath := filepath.Join(dir, outFilename)

	cmd := exec.Command("mkvmerge",
		"-o", outVideoPath,
		videoPath,
		"--language", "0:"+langCode,
		"--track-name", "0:BakaSub AI",
		"--default-track-flag", "0:yes",
		srtPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao fazer o merge com mkvmerge: %w | log: %s", err, string(output))
	}

	return outVideoPath, nil
}
