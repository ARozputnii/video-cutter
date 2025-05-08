package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CutRequest defines the input for cutting a video into time segments.
type CutRequest struct {
	Filename string `json:"filename"`
	Ranges   []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"ranges"`
}

// CutResponse defines the output after cutting and merging the video.
type CutResponse struct {
	Filename string `json:"filename"`
}

// CutServiceInterface defines the interface for video cut operations.
type CutServiceInterface interface {
	Cut(req CutRequest) (*CutResponse, error)
}

// FFmpegExecutor is an interface to abstract away the ffmpeg logic.
type FFmpegExecutor interface {
	CutSegment(inputPath, start, end, outFile string) error
	MergeSegments(listFile, outputPath string) error
}

// cutService manages cutting workflow using a processors.
type cutService struct {
	uploadDir string
	executor  FFmpegExecutor
}

// NewCutService creates a new services for cutting and merging videos.
func NewCutService(uploadDir string, executor FFmpegExecutor) CutServiceInterface {
	return &cutService{
		uploadDir: uploadDir,
		executor:  executor,
	}
}

func (s *cutService) Cut(req CutRequest) (*CutResponse, error) {
	originalName := filepath.Base(req.Filename)
	inputPath := filepath.Join(s.uploadDir, originalName)

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file not found")
	}

	tempDir, err := os.MkdirTemp("", "clips-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatalf("failed to remove temp dir: %v", err)
		}
	}(tempDir)

	var concatList bytes.Buffer
	for i, r := range req.Ranges {
		outFile := filepath.Join(tempDir, fmt.Sprintf("clip_%d.ts", i))
		if err := s.executor.CutSegment(inputPath, r.Start, r.End, outFile); err != nil {
			return nil, fmt.Errorf("cut segment error: %w", err)
		}
		concatList.WriteString("file '" + outFile + "'\n")
	}

	listPath := filepath.Join(tempDir, "inputs.txt")
	if err := os.WriteFile(listPath, concatList.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("failed to write concat list: %w", err)
	}

	outputPath := inputPath
	if err := s.executor.MergeSegments(listPath, outputPath); err != nil {
		return nil, fmt.Errorf("merge error: %w", err)
	}

	ext := filepath.Ext(originalName)
	_ = filepath.Walk(s.uploadDir, func(path string, info os.FileInfo, err error) error {
		if path != outputPath && strings.HasSuffix(path, ext) {
			_ = os.Remove(path)
		}
		return nil
	})

	return &CutResponse{Filename: req.Filename}, nil
}
