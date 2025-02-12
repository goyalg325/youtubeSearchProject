package youtube

import (
	"context"
	"fmt"
	"log"
	"time"
	"youtube-fetcher/internal/database"
	"youtube-fetcher/internal/models"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	apiKeys        []string
	currentKeyIndex int
}

func NewClient(apiKeys []string) *Client {
	return &Client{
		apiKeys:        apiKeys,
		currentKeyIndex: 0,
	}
}

func (c *Client) getNextAPIKey() string {
	c.currentKeyIndex = (c.currentKeyIndex + 1) % len(c.apiKeys)
	return c.apiKeys[c.currentKeyIndex]
}

func (c *Client) FetchAndStoreVideos(query string, db *database.DB) error {
    service, err := youtube.NewService(context.Background(), option.WithAPIKey(c.apiKeys[c.currentKeyIndex]))
    if err != nil {
        c.getNextAPIKey()
        return fmt.Errorf("error creating YouTube client: %v", err)
    }

    publishedAfter := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
    
     call := service.Search.List([]string{"id", "snippet"}).
        Q(query).
        Type("video").
        Order("date").
        PublishedAfter(publishedAfter).
        MaxResults(50).
        RegionCode("IN")

    response, err := call.Do()
    if err != nil {
        if isQuotaExceeded(err) {
            c.getNextAPIKey()
            return fmt.Errorf("quota exceeded, switched API key: %v", err)
        }
        return fmt.Errorf("error fetching videos: %v", err)
    }

    for _, item := range response.Items {
        publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
        
        video := &models.Video{
            ID:          item.Id.VideoId,
            Title:       item.Snippet.Title,
            Description: item.Snippet.Description,
            PublishedAt: publishedAt,
        }
        
        if item.Snippet.Thumbnails != nil {
            if item.Snippet.Thumbnails.Default != nil {
                video.ThumbnailDefault = item.Snippet.Thumbnails.Default.Url
            }
            if item.Snippet.Thumbnails.Medium != nil {
                video.ThumbnailMedium = item.Snippet.Thumbnails.Medium.Url
            }
            if item.Snippet.Thumbnails.High != nil {
                video.ThumbnailHigh = item.Snippet.Thumbnails.High.Url
            }
        }

        if err := db.StoreVideo(video); err != nil {
            log.Printf("Error storing video %s: %v", video.ID, err)
        }
    }

    return nil
}

func isQuotaExceeded(err error) bool {
	return err.Error() == "quotaExceeded"
}