package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"video-cutter/internal/handler"
)

const port = 3000

func main() {
	r := chi.NewRouter()

	r.Post("/upload", handler.UploadHandler)
	r.Post("/cut", handler.CutHandler)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("internal/frontend"))))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Starting server on port", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}
}
