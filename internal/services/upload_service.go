package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UploadServiceInterface defines the upload functionality contract.
type UploadServiceInterface interface {
	Upload(w http.ResponseWriter, r *http.Request) error
}

type uploadService struct {
	dir string
}

// NewUploadService creates a new upload service that stores files in the given directory.
func NewUploadService(uploadDir string) UploadServiceInterface {
	_ = os.MkdirAll(uploadDir, os.ModePerm)
	return &uploadService{dir: uploadDir}
}

func (us *uploadService) Upload(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		return fmt.Errorf("unable to parse form: %w", err)
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		return fmt.Errorf("unable to get file: %w", err)
	}
	defer file.Close()

	if err := us.clearUploadDir(); err != nil {
		return fmt.Errorf("failed to clear upload directory: %w", err)
	}

	filePath := filepath.Join(us.dir, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	_, _ = fmt.Fprint(w, header.Filename)
	return nil
}

// clearUploadDir removes all files from the upload directory.
func (us *uploadService) clearUploadDir() error {
	entries, err := os.ReadDir(us.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(us.dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
