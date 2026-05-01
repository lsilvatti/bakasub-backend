package services

import (
	"errors"
	"fmt"
	"os/exec"
)

const (
	toolFFmpeg     = "ffmpeg"
	toolMKVMerge   = "mkvmerge"
	toolMKVExtract = "mkvextract"
)

type ExternalToolStatus struct {
	Available bool   `json:"available"`
	Path      string `json:"path,omitempty"`
	Error     string `json:"error,omitempty"`
}

type VideoToolsStatus struct {
	FFmpeg                   ExternalToolStatus `json:"ffmpeg"`
	MKVMerge                 ExternalToolStatus `json:"mkvmerge"`
	MKVExtract               ExternalToolStatus `json:"mkvextract"`
	VideoProcessingAvailable bool               `json:"videoProcessingAvailable"`
	MissingTools             []string           `json:"missingTools,omitempty"`
}

type MissingVideoToolError struct {
	ToolName string
	Message  string
}

func (e *MissingVideoToolError) Error() string {
	return e.Message
}

func IsMissingVideoToolError(err error) bool {
	var missingToolErr *MissingVideoToolError
	return errors.As(err, &missingToolErr)
}

func CheckVideoTools() VideoToolsStatus {
	ffmpeg := checkExternalTool(toolFFmpeg)
	mkvmerge := checkExternalTool(toolMKVMerge)
	mkvextract := checkExternalTool(toolMKVExtract)

	missingTools := make([]string, 0, 3)
	if !ffmpeg.Available {
		missingTools = append(missingTools, toolFFmpeg)
	}
	if !mkvmerge.Available {
		missingTools = append(missingTools, toolMKVMerge)
	}
	if !mkvextract.Available {
		missingTools = append(missingTools, toolMKVExtract)
	}

	return VideoToolsStatus{
		FFmpeg:                   ffmpeg,
		MKVMerge:                 mkvmerge,
		MKVExtract:               mkvextract,
		VideoProcessingAvailable: len(missingTools) == 0,
		MissingTools:             missingTools,
	}
}

func resolveRequiredVideoTool(toolName string) (string, error) {
	status := checkExternalTool(toolName)
	if !status.Available {
		return "", &MissingVideoToolError{
			ToolName: toolName,
			Message:  status.Error,
		}
	}

	return status.Path, nil
}

func checkExternalTool(toolName string) ExternalToolStatus {
	toolPath, err := exec.LookPath(toolName)
	if err != nil {
		message := missingToolMessage(toolName)
		if !errors.Is(err, exec.ErrNotFound) {
			message = fmt.Sprintf("failed to locate %s in PATH: %v", toolName, err)
		}

		return ExternalToolStatus{
			Available: false,
			Error:     message,
		}
	}

	return ExternalToolStatus{
		Available: true,
		Path:      toolPath,
	}
}

func missingToolMessage(toolName string) string {
	switch toolName {
	case toolFFmpeg:
		return "ffmpeg is not installed or not available in PATH. Install FFmpeg and restart Bakasub."
	case toolMKVMerge, toolMKVExtract:
		return fmt.Sprintf("%s is not installed or not available in PATH. Install MKVToolNix and restart Bakasub.", toolName)
	default:
		return fmt.Sprintf("%s is not installed or not available in PATH.", toolName)
	}
}
