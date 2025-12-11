package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"snapcheck/internal/service"
)

// AnalyzeHandler handles the image analysis request.
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 21*1024*1024)
	if err := r.ParseMultipartForm(21 * 1024 * 1024); err != nil {
		http.Error(w, "File too large or invalid multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing 'file' field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var lastModified *time.Time
	if ts := r.FormValue("last_modified"); ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			lastModified = &t
		}
	}

	result, err := service.AnalyzeImage(data, lastModified)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
