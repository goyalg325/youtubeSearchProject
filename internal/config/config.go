package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    Port           string
    YouTubeAPIKeys []string
    SearchQuery    string
    FetchInterval  int
}

func New() *Config {
    fetchInterval, _ := strconv.Atoi(getEnv("FETCH_INTERVAL", "10"))
    
    return &Config{
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", ""),
        DBName:         getEnv("DB_NAME", "youtube_fetcher"),
        Port:           getEnv("PORT", "8080"),
        YouTubeAPIKeys: strings.Split(getEnv("YOUTUBE_API_KEYS", ""), ","),
        SearchQuery:    getEnv("SEARCH_QUERY", "golang"),
        FetchInterval:  fetchInterval,
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
