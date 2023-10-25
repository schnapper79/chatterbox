package chatterbox

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ProgressReader struct {
	reader       io.Reader
	total        int64
	bytesRead    int64
	progressFunc func(int64, int64)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.bytesRead += int64(n)
	pr.progressFunc(pr.bytesRead, pr.total)
	return n, err
}

func reportProgress(w http.ResponseWriter, bytesRead, total int64) {
	percentage := float64(bytesRead) / float64(total) * 100
	fmt.Fprintf(w, "data: {\"progress\": %.2f}\n\n", percentage)
}

func (s *Server) downloadModelHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Get the download URL from the query parameters
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to start download", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Initialize ProgressReader
	progressReader := &ProgressReader{
		reader: resp.Body,
		total:  resp.ContentLength,
		progressFunc: func(bytesRead, total int64) {
			reportProgress(w, bytesRead, total)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		},
	}

	// Create a new file to save the downloaded content
	out_split := strings.Split(url, "/")
	outFile, err := os.Create(s.ModelPath + "/" + out_split[len(out_split)-1])

	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Download the file (this will trigger the progress reporting)
	_, err = io.Copy(outFile, progressReader)
	if err != nil {
		http.Error(w, "Failed to complete download", http.StatusInternalServerError)
	}
}
