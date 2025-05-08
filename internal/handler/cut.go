package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type CutRequest struct {
	Filename       string `json:"filename"`
	DeleteOriginal bool   `json:"delete_original"`
	Ranges         []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"ranges"`
}

func CutHandler(w http.ResponseWriter, r *http.Request) {
	var req CutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	inputPath := filepath.Join("uploads", filepath.Base(req.Filename))
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		http.Error(w, "Input file not found", http.StatusNotFound)
		return
	}

	tempDir, err := os.MkdirTemp("", "clips-*")
	if err != nil {
		http.Error(w, "Cannot create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	var concatList bytes.Buffer
	for i, r := range req.Ranges {
		outFile := filepath.Join(tempDir, fmt.Sprintf("clip_%d.ts", i))
		cmd := exec.Command("ffmpeg",
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

	listPath := filepath.Join(tempDir, "inputs.txt")
	if err := ioutil.WriteFile(listPath, concatList.Bytes(), 0644); err != nil {
		http.Error(w, "Cannot write concat list", http.StatusInternalServerError)
		return
	}

	outputName := fmt.Sprintf("merged_%d.mp4", time.Now().Unix())
	outputPath := filepath.Join("uploads", outputName)

	mergeCmd := exec.Command("ffmpeg",
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

	if req.DeleteOriginal {
		_ = os.Remove(inputPath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Video saved",
		"filename": outputName,
	})
}
