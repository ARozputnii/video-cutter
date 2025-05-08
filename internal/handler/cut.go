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

	// Validate input path
	originalName := filepath.Base(req.Filename)
	inputPath := filepath.Join("uploads", originalName)

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		http.Error(w, "Input file not found", http.StatusNotFound)
		return
	}

	// Create temporary directory for .ts clips
	tempDir, err := os.MkdirTemp("", "clips-*")
	if err != nil {
		http.Error(w, "Cannot create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	var concatList bytes.Buffer

	// Cut video into temporary .ts segments using ffmpeg
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

		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("FFmpeg cut error: %v\nDetails: %s", err, stderr.String()), http.StatusInternalServerError)
			return
		}

		concatList.WriteString("file '" + outFile + "'\n")
	}

	// Write concat list file
	listPath := filepath.Join(tempDir, "inputs.txt")
	if err := ioutil.WriteFile(listPath, concatList.Bytes(), 0644); err != nil {
		http.Error(w, "Cannot write concat list", http.StatusInternalServerError)
		return
	}

	// Merge segments into original file (overwrite)
	outputPath := inputPath
	outputName := originalName

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

	// Remove all other video files with the same extension in uploads/
	ext := filepath.Ext(originalName)
	_ = filepath.Walk("uploads", func(path string, info fs.FileInfo, err error) error {
		if path != outputPath && strings.HasSuffix(path, ext) {
			_ = os.Remove(path)
		}
		return nil
	})

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Video saved",
		"filename": outputName,
	})
}

// getFFmpegBinary resolves the path to ffmpeg or ffmpeg.exe depending on OS and executable location
func getFFmpegBinary() string {
	name := "ffmpeg"

	if runtime.GOOS == "windows" {
		name = fmt.Sprintf("%s.exe", name)

		exePath, err := os.Executable()
		if err != nil {
			return name
		}
		return filepath.Join(filepath.Dir(exePath), name)
	}

	return name
}
