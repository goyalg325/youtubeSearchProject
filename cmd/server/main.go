package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"youtube-fetcher/internal/api"
	"youtube-fetcher/internal/config"
	"youtube-fetcher/internal/database"
	"youtube-fetcher/internal/youtube"

	"github.com/joho/godotenv"
)

func main() {
		if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	ytClient := youtube.NewClient(cfg.YouTubeAPIKeys)

	server := api.NewServer(cfg, db, ytClient)

	fetchTicker := time.NewTicker(time.Duration(cfg.FetchInterval) * time.Second)
	go func() {
		for range fetchTicker.C {
			if err := ytClient.FetchAndStoreVideos(cfg.SearchQuery, db); err != nil {
				log.Printf("Error fetching videos: %v", err)
			}
		}
	}()

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fetchTicker.Stop()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}