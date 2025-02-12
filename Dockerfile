# Dockerfile
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

RUN echo '#!/bin/bash\n\
if [ -z "$YOUTUBE_API_KEYS" ]; then\n\
    echo "Please enter your YouTube API keys (comma-separated for multiple keys):"\n\
    read api_keys\n\
    export YOUTUBE_API_KEYS=$api_keys\n\
fi\n\
\n\
echo "Starting application with configuration:"\n\
echo "YouTube API Keys: $YOUTUBE_API_KEYS"\n\
echo "Search Query: $SEARCH_QUERY"\n\
echo "Fetch Interval: $FETCH_INTERVAL seconds"\n\
\n\
exec ./main' > /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]