package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go/http3"
	"wlmwwx.duckdns.org/http3_upload/internal/config"
)

func setupUploadHandler(cfg *config.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Validate file size
		if header.Size > cfg.Server.MaxFileSize {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		// Create upload directory if it doesn't exist
		if err := os.MkdirAll(cfg.Server.UploadDir, 0755); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Create the file
		dst, err := os.Create(filepath.Join(cfg.Server.UploadDir, header.Filename))
		if err != nil {
			http.Error(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
	}
}

func setupDownloadHandler(cfg *config.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		filename := r.URL.Query().Get("file")
		if filename == "" {
			http.Error(w, "File parameter is required", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(cfg.Server.UploadDir, filename)

		// Prevent directory traversal
		if !filepath.HasPrefix(filePath, cfg.Server.UploadDir) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		http.ServeFile(w, r, filePath)
	}
}

func main() {
	cfg, err := config.LoadServerConfig("configs/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(cfg.Server.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", setupUploadHandler(cfg))
	mux.HandleFunc("/download", setupDownloadHandler(cfg))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http3.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting HTTP3 server on %s", addr)
	if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		log.Fatal(err)
	}
}