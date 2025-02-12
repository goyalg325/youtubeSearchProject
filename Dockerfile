# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

COPY <<EOF /app/docker-entrypoint.sh

if [ -z "\$YOUTUBE_API_KEYS" ]; then
    echo "Please enter your YouTube API keys (comma-separated for multiple keys):"
    read api_keys
    export YOUTUBE_API_KEYS=\$api_keys
fi

echo "Starting application with configuration:"
echo "YouTube API Keys: \$YOUTUBE_API_KEYS"
echo "Search Query: \$SEARCH_QUERY"
echo "Fetch Interval: \$FETCH_INTERVAL seconds"

exec ./main
EOF

RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]