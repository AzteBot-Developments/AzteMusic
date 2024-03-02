# Base Image
FROM golang:latest

WORKDIR /app

# Install project dependencies
RUN apt-get update && apt-get install -y iputils-ping
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of source code files
# TODO: Only copy required files
COPY . ./
COPY ./cmd/music-service/2/.env ./.env

# Build project
RUN CGO_ENABLED=0 GOOS=linux go build -o build/main ./internal/bot

# Other configs
VOLUME /app/logs

# Run the Discord bot service application
CMD sh -c "./build/main"