package services

import (
	"bakasub-backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
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

	utils.LogInfo("video", "info", "Scanning video for subtitle tracks", map[string]any{"videoPath": videoPath})
	utils.SendSSE("info", "video", "Scanning video for subtitle tracks...", map[string]any{"file": filepath.Base(videoPath)})

	probeCmd := exec.CommandContext(ctx, "mkvmerge", "-J", videoPath)
	var out bytes.Buffer
	probeCmd.Stdout = &out

	if err := probeCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			utils.LogError("video", "Scan tracks timeout", map[string]any{
				"videoPath": videoPath,
				"error":     "file read took too long",
			})
			utils.SendSSE("error", "video", "Scanning timeout exceeded.", nil)
			return nil, fmt.Errorf("timeout: file read took too long")
		}

		utils.LogError("video", "Scan tracks error", map[string]any{
			"videoPath": videoPath,
			"error":     err.Error(),
		})
		utils.SendSSE("error", "video", "Failed to scan video file.", nil)
		return nil, fmt.Errorf("error running probe on file: %w", err)
	}

	var info mkvMergeOutput
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		utils.LogError("video", "Failed to parse mkvmerge JSON", map[string]any{"error": err.Error()})
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

	utils.LogInfo("video", "success", "Scan tracks completed", map[string]any{
		"videoPath":  videoPath,
		"trackCount": len(subtitleTracks),
	})
	utils.SendSSE("success", "video", fmt.Sprintf("Found %d subtitle track(s).", len(subtitleTracks)), nil)

	return subtitleTracks, nil
}

func (s *VideoService) ExtractSubtitle(videoPath string, subtitleId int) (string, error) {
	tracks, err := s.ScanSubtitles(videoPath)
	if err != nil || len(tracks) == 0 {
		utils.LogError("video", "Extract subtitle error: no tracks found", map[string]any{
			"videoPath":  videoPath,
			"subtitleId": subtitleId,
			"error":      err.Error(),
		})
		utils.SendSSE("error", "video", "Failed to find subtitle tracks for extraction.", nil)
		return "", fmt.Errorf("failed to read video tracks or no tracks found: %v", err)
	}

	utils.LogInfo("video", "info", "Subtitle extraction initialized", map[string]any{
		"videoPath":  videoPath,
		"subtitleId": subtitleId,
	})

	lang := "und"
	codec := ""
	trackFound := false
	for _, t := range tracks {
		if t.ID == subtitleId {
			lang = t.Language
			codec = t.Codec
			trackFound = true
			break
		}
	}

	if !trackFound {
		utils.LogError("video", "Extract subtitle error: track not found", map[string]any{
			"videoPath":  videoPath,
			"subtitleId": subtitleId,
		})
		utils.SendSSE("error", "video", "Selected subtitle track not found in file.", nil)
		return "", fmt.Errorf("track with ID %d not found", subtitleId)
	}

	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	outExt := ".srt"
	upperCodec := strings.ToUpper(codec)
	if strings.Contains(upperCodec, "ASS") || strings.Contains(upperCodec, "SSA") {
		outExt = ".ass"
	}

	subFilename := fmt.Sprintf("%s_%d_%s%s", nameWithoutExt, subtitleId, lang, outExt)
	subPath := filepath.Join(dir, subFilename)

	fileExt := strings.ToLower(ext)

	utils.SendSSE("progress", "video", "Extracting subtitle from video...", map[string]any{"track_id": subtitleId})

	if fileExt == ".mkv" {
		err = s.ExtractWithMKVToolnix(videoPath, subtitleId, subPath)
	} else {
		err = s.ExtractWithFFmpeg(videoPath, subtitleId, subPath)
	}

	if err != nil {
		return "", err
	}

	content, errRead := os.ReadFile(subPath)
	if errRead == nil && len(content) > 0 {
		header := string(content)
		if len(header) > 100 {
			header = header[:100]
		}

		if strings.Contains(header, "[Script Info]") && strings.HasSuffix(subPath, ".srt") {
			newPath := strings.TrimSuffix(subPath, ".srt") + ".ass"
			os.Rename(subPath, newPath)
			subPath = newPath
		}
	}

	utils.SendSSE("success", "video", "Subtitle extracted successfully!", map[string]any{"output": filepath.Base(subPath)})
	return subPath, nil
}

func (s *VideoService) ExtractWithMKVToolnix(videoPath string, subtitleId int, subPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	utils.LogInfo("video", "info", "Starting MKVToolnix extraction", map[string]any{"tool": "mkvextract", "video": videoPath})

	extractCmd := exec.CommandContext(ctx, "mkvextract", videoPath, "tracks", fmt.Sprintf("%d:%s", subtitleId, subPath))
	output, err := extractCmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		utils.LogError("video", "MKVToolnix extraction timeout", map[string]any{"video": videoPath})
		utils.SendSSE("error", "video", "Subtitle extraction timed out.", nil)
		return fmt.Errorf("timeout: mkvextract extraction exceeded 10 minute limit")
	}
	if err != nil {
		utils.LogError("video", "MKVToolnix extraction failed", map[string]any{
			"error":  err.Error(),
			"output": string(output),
		})
		utils.SendSSE("error", "video", "Failed to extract subtitle using MKVToolnix.", nil)
		return fmt.Errorf("error running mkvextract: %w | log: %s", err, string(output))
	}

	utils.LogInfo("video", "success", "MKVToolnix extraction successful", map[string]any{"output": subPath})
	return nil
}

func (s *VideoService) ExtractWithFFmpeg(videoPath string, subtitleId int, subPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	utils.LogInfo("video", "info", "Starting FFmpeg extraction", map[string]any{"tool": "ffmpeg", "video": videoPath})

	extractCmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", videoPath, "-map", fmt.Sprintf("0:%d", subtitleId), subPath)
	output, err := extractCmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		utils.LogError("video", "FFmpeg extraction timeout", map[string]any{"video": videoPath})
		utils.SendSSE("error", "video", "Subtitle extraction timed out.", nil)
		return fmt.Errorf("timeout: ffmpeg extraction exceeded 10 minute limit")
	}
	if err != nil {
		utils.LogError("video", "FFmpeg extraction failed", map[string]any{
			"error":  err.Error(),
			"output": string(output),
		})
		utils.SendSSE("error", "video", "Failed to extract subtitle using FFmpeg.", nil)
		return fmt.Errorf("error extracting with FFmpeg: %w | log: %s", err, string(output))
	}

	utils.LogInfo("video", "success", "FFmpeg extraction successful", map[string]any{"output": subPath})
	return nil
}

func (s *VideoService) MergeSubtitle(videoPath string, subPath string, langCode string, timeoutMinutes int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

	dir := filepath.Dir(videoPath)
	base := filepath.Base(videoPath)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	outFilename := fmt.Sprintf("%s_bakasub_%s.mkv", nameWithoutExt, langCode)
	outVideoPath := filepath.Join(dir, outFilename)

	utils.LogInfo("video", "info", "Initializing subtitle merge process", map[string]any{
		"videoPath": videoPath,
		"subPath":   subPath,
		"lang":      langCode,
	})
	utils.SendSSE("progress", "video", "Muxing translated subtitle into new MKV video...", map[string]any{"file": outFilename})

	cmd := exec.CommandContext(ctx, "mkvmerge",
		"-o", outVideoPath,
		videoPath,
		"--language", "0:"+langCode,
		"--track-name", "0:BakaSub AI",
		"--default-track-flag", "0:yes",
		subPath,
	)

	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		utils.LogError("video", "Merge process timeout", map[string]any{"video": videoPath})
		utils.SendSSE("error", "video", "Video merge process timed out.", nil)
		return "", fmt.Errorf("timeout: merge process exceeded limit")
	}
	if err != nil {
		utils.LogError("video", "Failed to merge with mkvmerge", map[string]any{
			"error":  err.Error(),
			"output": string(output),
		})
		utils.SendSSE("error", "video", "Failed to merge translated subtitle into video.", nil)
		return "", fmt.Errorf("failed to merge with mkvmerge: %w | log: %s", err, string(output))
	}

	utils.LogInfo("video", "success", "Merge process completed successfully", map[string]any{"outVideoPath": outVideoPath})
	utils.SendSSE("success", "video", "Video processing complete!", map[string]any{"output": outFilename})

	return outVideoPath, nil
}
