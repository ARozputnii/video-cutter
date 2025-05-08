package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CutRequest struct {
	Filename string `json:"filename"`
	Ranges   []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"ranges"`
}

func CutHandler(w http.ResponseWriter, r *http.Request) {
	var req CutRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Clean and validate input filename
	originalName := filepath.Base(req.Filename)
	inputPath := filepath.Join("uploads", originalName)

	// Check if file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		http.Error(w, "Input file not found", http.StatusNotFound)
		return
	}

	// Create a temporary directory for intermediate clips
	tempDir, err := os.MkdirTemp("", "clips-*")
	if err != nil {
		http.Error(w, "Cannot create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir) // Auto-delete when function exits

	var concatList bytes.Buffer

	// For each range, create a temporary .ts clip
	for i, r := range req.Ranges {
		outFile := filepath.Join(tempDir, fmt.Sprintf("clip_%d.ts", i))
		cmd := exec.Command(getFFmpegBinary(),
			"-i", inputPath,
			"-ss", r.Start,
			"-to", r.End,
			"-c", "copy",
			"-bsf:v", "h264_mp4toannexb",
			"-f", "mpegts",
			outFile,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		// Run ffmpeg for this range
		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("FFmpeg cut error: %v\nDetails: %s", err, stderr.String()), http.StatusInternalServerError)
			return
		}

		concatList.WriteString("file '" + outFile + "'\n")
	}

	// Save ffmpeg concat list
	listPath := filepath.Join(tempDir, "inputs.txt")
	if err := ioutil.WriteFile(listPath, concatList.Bytes(), 0644); err != nil {
		http.Error(w, "Cannot write concat list", http.StatusInternalServerError)
		return
	}

	// Overwrite original file with merged result
	outputPath := inputPath
	outputName := originalName

	// Merge all .ts clips into one .mp4
	mergeCmd := exec.Command(getFFmpegBinary(),
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", listPath,
		"-c", "copy",
		outputPath,
	)

	var mergeErr bytes.Buffer
	mergeCmd.Stderr = &mergeErr

	if err := mergeCmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("FFmpeg merge error: %v\nDetails: %s", err, mergeErr.String()), http.StatusInternalServerError)
		return
	}

	// Clean up all other video files in uploads/ except the one we just saved
	ext := filepath.Ext(originalName)
	_ = filepath.Walk("uploads", func(path string, info fs.FileInfo, err error) error {
		if path != outputPath && strings.HasSuffix(path, ext) {
			os.Remove(path)
		}
		return nil
	})

	// Send response with filename
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Video saved",
		"filename": outputName,
	})
}

// getFFmpegBinary returns platform-specific binary name
func getFFmpegBinary() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}
