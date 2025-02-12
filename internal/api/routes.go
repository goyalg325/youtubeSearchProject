package api

import (
    "context"
    "net/http"
    "sync"
    "time"
    "youtube-fetcher/internal/config"
    "youtube-fetcher/internal/database"
    "youtube-fetcher/internal/youtube"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

type rateLimiter struct {
    sync.Mutex
    requests map[string][]time.Time
}

func newRateLimiter() *rateLimiter {
    return &rateLimiter{
        requests: make(map[string][]time.Time),
    }
}

func (rl *rateLimiter) Allow(ip string) bool {
    rl.Lock()
    defer rl.Unlock()

    now := time.Now()
    window := now.Add(-time.Minute)

    if _, exists := rl.requests[ip]; !exists {
        rl.requests[ip] = []time.Time{}
    }

    requests := rl.requests[ip]
    valid := requests[:0]
    for _, t := range requests {
        if t.After(window) {
            valid = append(valid, t)
        }
    }
    rl.requests[ip] = valid

    if len(valid) >= 100 {
        return false
    }

    rl.requests[ip] = append(rl.requests[ip], now)
    return true
}

type Server struct {
    cfg         *config.Config
    db          *database.DB
    youtube     *youtube.Client
    router      *chi.Mux
    rateLimiter *rateLimiter
}

func NewServer(cfg *config.Config, db *database.DB, yt *youtube.Client) *Server {
    s := &Server{
        cfg:         cfg,
        db:          db,
        youtube:     yt,
        router:      chi.NewRouter(),
        rateLimiter: newRateLimiter(),
    }

    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    s.router.Use(middleware.Logger)
    s.router.Use(middleware.Recoverer)
    s.router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    s.router.Get("/health", s.handleHealthCheck)
    s.router.Route("/api", func(r chi.Router) {
        r.Get("/videos", s.handleGetVideos)
    })
}

func (s *Server) Start() error {
    return http.ListenAndServe(":"+s.cfg.Port, s.router)
}

func (s *Server) Shutdown(ctx context.Context) error {
    return nil
}