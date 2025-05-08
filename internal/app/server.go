package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"video-cutter/internal/handlers"
	"video-cutter/internal/processors"
	"video-cutter/internal/services"
)

// ServerInterface defines the methods for starting and stopping the server.
type ServerInterface interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// appServer implements ServerInterface and holds internal HTTP server.
type appServer struct {
	httpServer *http.Server
}

// NewServer creates and configures a new appServer instance.
func NewServer(cfg Config) ServerInterface {
	uploadSvc := services.NewUploadService(cfg.UploadDir)
	cutExecutor := processors.NewFFmpegProcessor()
	cutSvc := services.NewCutService(cfg.UploadDir, cutExecutor)

	uploadHandler := handlers.NewUploadHandler(uploadSvc)
	cutHandler := handlers.NewCutHandler(cutSvc)

	r := chi.NewRouter()
	r.Post("/upload", uploadHandler.Upload)
	r.Post("/cut", cutHandler.Cut)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(cfg.FrontendDir))))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: r,
	}

	return &appServer{httpServer: srv}
}

// Start runs the HTTP server and blocks until interrupted.
func (a *appServer) Start() error {
	go func() {
		log.Println("Server starting on", a.httpServer.Addr)
		log.Println("App available at: http://localhost" + a.httpServer.Addr + "/")
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Received shutdown signal")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.Shutdown(ctx)
}

// Shutdown gracefully stops the HTTP server.
func (a *appServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server stopped gracefully.")
	return nil
}
