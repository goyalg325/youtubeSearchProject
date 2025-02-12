package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

func (s *Server) handleGetVideos(w http.ResponseWriter, r *http.Request) {
    // Rate limiting check
    clientIP := r.RemoteAddr
    if !s.rateLimiter.Allow(clientIP) {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }

    perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
    if perPage < 1 {
        perPage = 10
    }

    sortDir := r.URL.Query().Get("sort")
    if sortDir != "asc" {
        sortDir = "desc"
    }

    videos, total, err := s.db.GetVideos(page, perPage, sortDir)
    if err != nil {
        log.Printf("Error in handleGetVideos: %v", err)
        http.Error(w, fmt.Sprintf("Error fetching videos: %v", err), http.StatusInternalServerError)
        return
    }

    response := struct {
        Videos     interface{} `json:"videos"`
        Pagination struct {
            CurrentPage int `json:"current_page"`
            TotalPages  int `json:"total_pages"`
            TotalItems  int `json:"total_items"`
        } `json:"pagination"`
    }{
        Videos: videos,
        Pagination: struct {
            CurrentPage int `json:"current_page"`
            TotalPages  int `json:"total_pages"`
            TotalItems  int `json:"total_items"`
        }{
            CurrentPage: page,
            TotalPages:  (total + perPage - 1) / perPage,
            TotalItems:  total,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}