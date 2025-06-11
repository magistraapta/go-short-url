package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"short-url/dto"
)

var storage = map[string]string{}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/shorten", shortenHandler)
	router.HandleFunc("/", redirectHandler)
	router.HandleFunc("/check", checkStorage)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func generateShortURL() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func shortUrl(originalURL string) string {
	shortURL := generateShortURL()
	storage[shortURL] = originalURL
	return shortURL
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "Method not allowed"})
		return
	}

	var req dto.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if !isValidURL(req.URL) {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid URL format"})
		return
	}

	shortURL := shortUrl(req.URL)
	response := dto.ShortenResponse{
		ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortURL),
		OriginalURL: req.URL,
	}

	writeJSON(w, http.StatusOK, response)
}

func getShortUrl(shortURL string) string {
	if url, ok := storage[shortURL]; ok {
		return url
	}
	return ""
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	if shortURL == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "Short URL is required"})
		return
	}

	originalURL := getShortUrl(shortURL)
	if originalURL == "" {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: "URL not found"})
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func checkStorage(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, storage)
}
