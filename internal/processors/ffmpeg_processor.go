package processors

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// FFmpegProcessor executes actual ffmpeg commands.
type FFmpegProcessor struct{}

// NewFFmpegProcessor returns a new ffmpeg executor.
func NewFFmpegProcessor() *FFmpegProcessor {
	return &FFmpegProcessor{}
}

// CutSegment cuts a portion of a video using ffmpeg and saves it to a temporary .ts file.
func (f *FFmpegProcessor) CutSegment(inputPath, start, end, outFile string) error {
	cmd := exec.Command(getFFmpegBinary(),
		"-i", inputPath,
		"-ss", start,
		"-to", end,
		"-c", "copy",
		"-f", "mpegts",
		outFile,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg cut error: %v: %s", err, stderr.String())
	}
	return nil
}

// MergeSegments merges multiple .ts segments into a single output file.
func (f *FFmpegProcessor) MergeSegments(listFile, outputPath string) error {
	cmd := exec.Command(getFFmpegBinary(),
		"-y", "-f", "concat", "-safe", "0",
		"-i", listFile,
		"-c", "copy",
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg merge error: %v: %s", err, stderr.String())
	}
	return nil
}

// getFFmpegBinary returns the ffmpeg binary name depending on the OS.
func getFFmpegBinary() string {
	name := "ffmpeg"
	if runtime.GOOS == "windows" {
		name += ".exe"
		exePath, err := os.Executable()
		if err != nil {
			return name
		}
		return filepath.Join(filepath.Dir(exePath), name)
	}
	return name
}
