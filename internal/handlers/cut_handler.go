package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"video-cutter/internal/services"
)

// CutHandler handles HTTP requests for cutting video files.
type CutHandler struct {
	service services.CutServiceInterface
}

// NewCutHandler returns a new instance of CutHandler.
func NewCutHandler(s services.CutServiceInterface) *CutHandler {
	return &CutHandler{service: s}
}

// Cut handles the POST /cut request to process and trim a video file.
func (h *CutHandler) Cut(w http.ResponseWriter, r *http.Request) {
	var req services.CutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Cut(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Fatalf("Failed to write response: %v", err)
		return
	}
}
