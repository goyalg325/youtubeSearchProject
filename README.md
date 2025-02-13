# YouTube Video Fetcher

A Go service that continuously fetches YouTube videos based on specified search queries and stores them in PostgreSQL. Features include API key rotation, rate limiting, pagination, and automatic cleanup of old records.

## Features

- **YouTube Data Integration**
  - Fetches videos using YouTube Data API v3
  - Multiple API key support with automatic rotation
  - Region-specific search (currently set to India)
  - Configurable search queries

- **Data Management**
  - PostgreSQL storage with efficient indexing
  - Automatic cleanup of videos older than 24 hours
  - Upsert support to prevent duplicates

- **API Features**
  - Rate limiting (100 requests per minute per IP)
  - Pagination support
  - Sorting by publication date (asc/desc)
  - CORS enabled

## Requirements

- Go 1.21 or higher
- PostgreSQL 14 or higher
- YouTube Data API key(s)

## Getting Started

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd youtube-fetcher
   ```

2. **Set Up Environment**
   ```bash
   cp .env.example .env
   ```

3. **Database Setup**
   ```sql
   CREATE DATABASE youtube_fetcher;
   ```

4. **Install Dependencies**
   ```bash
   go mod download
   ```

5. **Run the Application**
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Get Videos
```http
GET /api/videos
```

#### Query Parameters
| Parameter | Default | Description |
|-----------|---------|-------------|
| page      | 1       | Page number |
| per_page  | 10      | Items per page |
| sort      | "desc"  | Sort by date ["asc", "desc"] |

#### Response
```json
{
  "videos": [
    {
      "id": "video_id",
      "title": "Video Title",
      "description": "Video Description",
      "published_at": "2024-02-12T19:45:26Z",
      "thumbnail_default": "url",
      "thumbnail_medium": "url",
      "thumbnail_high": "url",
      "created_at": "2024-02-12T19:45:26Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 100
  }
}
```
### Health Check
```http
GET /health
```

#### Response
```json
{
  "status": "ok"
}
```

## Technical Details

### Rate Limiting
- **Limit**: 100 requests per minute per IP
- **Implementation**: Thread-safe with mutex lock
- **Response Code**: 429 Too Many Requests (when limit exceeded)

### Video Management
- **Retention Period**: 24 hours
- **Cleanup**: Automatic removal of older videos
- **Storage**: PostgreSQL with optimized indexing
- **Deduplication**: Upsert support for existing videos

### Configuration
- **Fetch Interval**: Configurable (default: 60 seconds)
- **API Keys**: Multiple key support with rotation
- **Search Terms**: Pipe-separated query strings
- **Region**: Configurable (default: IN)

## Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=youtube_fetcher

# YouTube API Configuration
YOUTUBE_API_KEYS=key1,key2,key3

# Search Configuration
SEARCH_QUERY=breaking news|live news|news today

# Application Configuration
FETCH_INTERVAL=60  # seconds
PORT=8080
```

## Project Structure
```plaintext
youtube_fetcher/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── routes.go
│   │   └── handlers.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── db.go
│   ├── models/
│   │   └── video.go
│   └── youtube/
│       └── client.go
└── go.mod
└── go.sum
```

## Implementation Details

### Database Schema
```sql
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
```

### Key Features Implementation
- **Rate Limiting**: Thread-safe implementation using mutex
- **Database**: Using `sqlx` for enhanced database operations
- **API**: Chi router with middleware support
- **CORS**: Enabled for cross-origin requests
- **Error Handling**: Graceful error management with proper status codes
## Error Handling

### API Errors
| Status Code | Description | Handling |
|-------------|-------------|----------|
| 429 | Rate Limit Exceeded | Wait for rate limit window reset |
| 500 | Internal Server Error | Check server logs for details |
| 400 | Bad Request | Verify request parameters |

### YouTube API
- Automatic key rotation on quota exceeded
- Graceful handling of API limitations
- Configurable retry mechanisms

### Database
- Connection pool management
- Transaction handling
- Error logging and recovery

## API Response Examples

### Successful Response
```json
{
  "videos": [
    {
      "id": "video_id",
      "title": "Video Title",
      "description": "Video Description",
      "published_at": "2024-02-12T19:45:26Z",
      "thumbnail_default": "url",
      "thumbnail_medium": "url",
      "thumbnail_high": "url",
      "created_at": "2024-02-12T19:45:26Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 100
  }
}
```

### Error Response
```json
{
  "error": "Rate limit exceeded",
  "status": 429,
  "message": "Please try again later"
}
```

## Performance Considerations

### Database Optimization
- Indexed queries for faster retrieval
- Connection pooling
- Prepared statements

### API Efficiency
- Rate limiting per IP
- Response caching (where applicable)
- Efficient error handling

## Development Guide

### Prerequisites
```bash
# Install Go
go version  # Should be 1.21 or higher

# Install PostgreSQL
psql --version  # Should be 14 or higher
```

### Local Development
1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Update environment variables in `.env`

3. Create database:
   ```sql
   CREATE DATABASE youtube_fetcher;
   ```

4. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

<img width="637" alt="apiworking" src="https://github.com/user-attachments/assets/93e18489-1e5c-40f8-8a44-2566400357d1" />

<img width="629" alt="apiworking1" src="https://github.com/user-attachments/assets/3c60f414-daff-4ada-8197-a589559a8c5d" />

<img width="636" alt="apiworking2" src="https://github.com/user-attachments/assets/d5bacc43-54e4-4ee3-8eee-31b372bf4796" />




