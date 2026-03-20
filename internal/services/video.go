package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

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

type VideoService struct{}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (s *VideoService) ScanSubtitles(videoPath string) ([]SubtitleTrack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	probeCmd := exec.CommandContext(ctx, "mkvmerge", "-J", videoPath)
	var out bytes.Buffer
	probeCmd.Stdout = &out

	if err := probeCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout: file read took too long")
		}
		return nil, fmt.Errorf("error running probe on file: %w", err)
	}

	var info mkvMergeOutput
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return nil, fmt.Errorf("error parsing mkvmerge JSON: %w", err)
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

func (s *VideoService) ExtractSubtitle(videoPath string, subtitleId int) (string, error) {
	tracks, err := s.ScanSubtitles(videoPath)
	if err != nil || len(tracks) == 0 {
		return "", fmt.Errorf("failed to read video tracks or no tracks found: %v", err)
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
		return "", fmt.Errorf("track with ID %d not found", subtitleId)
	}

	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	srtFilename := fmt.Sprintf("%s_%d_%s.srt", nameWithoutExt, subtitleId, lang)
	srtPath := filepath.Join(dir, srtFilename)

	fileExt := strings.ToLower(ext)

	if fileExt == ".mkv" {
		err = s.ExtractWithMKVToolnix(videoPath, subtitleId, srtPath)
	} else {
		err = s.ExtractWithFFmpeg(videoPath, subtitleId, srtPath)
	}

	if err != nil {
		return "", err
	}

	return srtPath, nil
}

func (s *VideoService) ExtractWithMKVToolnix(videoPath string, subtitleId int, srtPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	extractCmd := exec.CommandContext(ctx, "mkvextract", videoPath, "tracks", fmt.Sprintf("%d:%s", subtitleId, srtPath))
	output, err := extractCmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timeout: mkvextract extraction exceeded 10 minute limit")
	}
	if err != nil {
		return fmt.Errorf("error running mkvextract: %w | log: %s", err, string(output))
	}
	return nil
}

func (s *VideoService) ExtractWithFFmpeg(videoPath string, subtitleId int, srtPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	extractCmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", videoPath, "-map", fmt.Sprintf("0:%d", subtitleId), srtPath)
	output, err := extractCmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timeout: ffmpeg extraction exceeded 10 minute limit")
	}
	if err != nil {
		return fmt.Errorf("error extracting with FFmpeg: %w | log: %s", err, string(output))
	}
	return nil
}

func (s *VideoService) MergeSubtitle(videoPath string, srtPath string, langCode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	outFilename := fmt.Sprintf("%s_bakasub_%s.mkv", nameWithoutExt, langCode)
	outVideoPath := filepath.Join(dir, outFilename)

	cmd := exec.CommandContext(ctx, "mkvmerge",
		"-o", outVideoPath,
		videoPath,
		"--language", "0:"+langCode,
		"--track-name", "0:BakaSub AI",
		"--default-track-flag", "0:yes",
		srtPath,
	)

	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout: merge process exceeded 15 minute limit")
	}
	if err != nil {
		return "", fmt.Errorf("failed to merge with mkvmerge: %w | log: %s", err, string(output))
	}

	return outVideoPath, nil
}
