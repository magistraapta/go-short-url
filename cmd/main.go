package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"short-url/config"
	"short-url/dto"
	"short-url/models"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	config.LoadEnv()
	db = config.ConnectDatabase()

	// Auto migrate the schema
	db.AutoMigrate(&models.URL{})

	router := http.NewServeMux()

	router.HandleFunc("/shorten", shortenHandler)
	router.HandleFunc("/", redirectHandler)
	router.HandleFunc("/check", checkStorage)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
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
	if shortURL == "" {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create short URL"})
		return
	}

	response := dto.ShortenResponse{
		ShortURL:    fmt.Sprintf("http://localhost:8080/%s", shortURL),
		OriginalURL: req.URL,
	}

	writeJSON(w, http.StatusOK, response)
}

func getShortUrl(shortURL string) string {
	var urlRecord models.URL
	if err := db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
		return ""
	}
	return urlRecord.OriginalURL
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
	var urls []models.URL
	if err := db.Find(&urls).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch URLs"})
		return
	}
	writeJSON(w, http.StatusOK, urls)
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

	// Create new URL record in database
	urlRecord := models.URL{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	if err := db.Create(&urlRecord).Error; err != nil {
		log.Printf("Error creating URL record: %v", err)
		return ""
	}

	return shortURL
}
