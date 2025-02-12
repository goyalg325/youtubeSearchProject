package database

import (
	"fmt"
	"youtube-fetcher/internal/config"
	"youtube-fetcher/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func New(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS videos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		published_at TIMESTAMP NOT NULL,
		thumbnail_default TEXT,
		thumbnail_medium TEXT,
		thumbnail_high TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_published_at ON videos(published_at);
	`
	_, err := db.Exec(schema)
	return err
}

func (db *DB) StoreVideo(video *models.Video) error {
    query := `
        INSERT INTO videos (
            id, title, description, published_at, 
            thumbnail_default, thumbnail_medium, thumbnail_high
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE
        SET 
            title = $2,
            description = $3,
            published_at = $4,
            thumbnail_default = $5,
            thumbnail_medium = $6,
            thumbnail_high = $7
    `
    _, err := db.Exec(query,
        video.ID,
        video.Title,
        video.Description,
        video.PublishedAt,
        video.ThumbnailDefault,
        video.ThumbnailMedium,
        video.ThumbnailHigh,
    )
    return err
}

func (db *DB) GetVideos(page, perPage int, sortDir string) ([]models.Video, int, error) {
    var total int
    if err := db.Get(&total, "SELECT COUNT(*) FROM videos"); err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * perPage
    orderBy := "DESC"
    if sortDir == "asc" {
        orderBy = "ASC"
    }

    query := fmt.Sprintf(`
        SELECT 
            id, title, description, published_at,
            thumbnail_default, thumbnail_medium, thumbnail_high,
            created_at
        FROM videos
        ORDER BY published_at %s
        LIMIT $1 OFFSET $2
    `, orderBy)

    var videos []models.Video
    if err := db.Select(&videos, query, perPage, offset); err != nil {
        return nil, 0, fmt.Errorf("error selecting videos: %v", err)
    }

    return videos, total, nil
}