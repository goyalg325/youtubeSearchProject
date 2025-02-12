FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/server

RUN echo '#!/bin/bash\n\
if [ -z "$YOUTUBE_API_KEYS" ]; then\n\
    echo "Enter your YouTube API keys (comma-separated for multiple keys):"\n\
    read -p "> " api_keys\n\
    export YOUTUBE_API_KEYS=$api_keys\n\
fi\n\
./main' > /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]