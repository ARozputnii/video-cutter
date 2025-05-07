package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CutRequest struct {
	Filename   string `json:"filename"`
	Start      string `json:"start"`
	End        string `json:"end"`
	DeleteOrig bool   `json:"delete_original"`
}

func CutHandler(w http.ResponseWriter, r *http.Request) {
	var req CutRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	inputPath := filepath.Join("uploads", req.Filename)
	outputFile := strings.TrimSuffix(req.Filename, filepath.Ext(req.Filename)) + "_cut.mp4"
	outputPath := filepath.Join("uploads", outputFile)

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-ss", req.Start, "-to", req.End, "-c", "copy", "-avoid_negative_ts", "1", outputPath)
	err := cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("FFmpeg error: %v", err), http.StatusInternalServerError)
		return
	}

	if req.DeleteOrig {
		if err := os.Remove(inputPath); err != nil {
			http.Error(w, fmt.Sprintf("Clip saved as %s, but failed to delete original: %v", outputFile, err), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, outputFile)
}
