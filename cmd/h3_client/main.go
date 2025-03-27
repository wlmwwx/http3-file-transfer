package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/quic-go/quic-go/http3"
	"wlmwwx.duckdns.org/http3_upload/internal/config"
)

func uploadFile(client *http.Client, cfg *config.ClientConfig, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	url := fmt.Sprintf("https://%s:%d/upload", cfg.Client.ServerHost, cfg.Client.ServerPort)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("File %s uploaded successfully\n", filepath.Base(filePath))
	return nil
}

func downloadFile(client *http.Client, cfg *config.ClientConfig, filename string) error {
	url := fmt.Sprintf("https://%s:%d/download?file=%s", cfg.Client.ServerHost, cfg.Client.ServerPort, filename)

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := os.MkdirAll(cfg.Client.DownloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %v", err)
	}

	filePath := filepath.Join(cfg.Client.DownloadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("File downloaded successfully to %s\n", filePath)
	return nil
}

func main() {
	cfg, err := config.LoadClientConfig("configs/client.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Command line flags
	upload := flag.String("upload", "", "File to upload")
	download := flag.String("download", "", "File to download")
	flag.Parse()

	if *upload == "" && *download == "" {
		log.Fatal("Please specify either -upload or -download flag")
	}

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Client.SkipVerify,
		},
	}
	defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
		Timeout:   time.Duration(30) * time.Second,
	}

	if *upload != "" {
		if err := uploadFile(client, cfg, *upload); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
	}

	if *download != "" {
		if err := downloadFile(client, cfg, *download); err != nil {
			log.Fatalf("Download failed: %v", err)
		}
	}
}