package node

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Add new handler function
func handleFetchLocalUrl(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the URL from request body
	var requestBody struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL is localhost/127.0.0.1
	if !isLocalURL(requestBody.URL) {
		http.Error(w, "Only localhost URLs are allowed", http.StatusForbidden)
		return
	}

	// Make the request
	resp, err := http.Get(requestBody.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set response status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		// Note: Can't write error to response here as headers are already sent
		log.Printf("Error copying response: %v", err)
	}
}

// Helper function to validate localhost URLs
func isLocalURL(urlStr string) bool {
	u, err := http.ParseURL(urlStr)
	if err != nil {
		return false
	}

	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1"
}
