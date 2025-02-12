package models

import "time"

type Video struct {
    ID              string    `json:"id" db:"id"`
    Title           string    `json:"title" db:"title"`
    Description     string    `json:"description" db:"description"`
    PublishedAt     time.Time `json:"published_at" db:"published_at"`
    ThumbnailDefault string   `json:"thumbnail_default" db:"thumbnail_default"`
    ThumbnailMedium  string   `json:"thumbnail_medium" db:"thumbnail_medium"`
    ThumbnailHigh    string   `json:"thumbnail_high" db:"thumbnail_high"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type Pagination struct {
    CurrentPage int `json:"current_page"`
    TotalPages  int `json:"total_pages"`
    TotalItems  int `json:"total_items"`
}

type PaginatedResponse struct {
    Videos     []Video    `json:"videos"`
    Pagination Pagination `json:"pagination"`
}