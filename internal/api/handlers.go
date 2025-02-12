package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "youtube-fetcher/internal/models"
)

func (s *Server) handleGetVideos(w http.ResponseWriter, r *http.Request) {
    if err := s.db.Ping(); err != nil {
        log.Printf("Database connection error: %v", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
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

    if videos == nil {
        videos = []models.Video{} 
    }

    totalPages := (total + perPage - 1) / perPage

    response := models.PaginatedResponse{
        Videos: videos,
        Pagination: models.Pagination{
            CurrentPage: page,
            TotalPages:  totalPages,
            TotalItems:  total,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}