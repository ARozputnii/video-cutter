package handlers

import (
	"fmt"
	"net/http"
	"video-cutter/internal/services"
)

// UploadHandler handles HTTP requests for uploading video files.
type UploadHandler struct {
	service services.UploadServiceInterface
}

// NewUploadHandler returns a new instance of UploadHandler.
func NewUploadHandler(s services.UploadServiceInterface) *UploadHandler {
	return &UploadHandler{service: s}
}

// Upload handles the POST /upload request to store an uploaded video file.
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := h.service.Upload(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
	}
}
