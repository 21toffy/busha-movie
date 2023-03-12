# Build stage
FROM golang:alpine AS builder

# Set the working directory to the app directory
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the application binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Final stage
FROM alpine:3.14

# Set the working directory to the app directory
WORKDIR /app

# Copy the application binary from the build stage
COPY --from=builder /app/main .

COPY wait-for.sh /app/wait-for.sh
COPY start.sh /app/start.sh

RUN chmod +x /app/wait-for.sh
RUN chmod +x /app/start.sh


COPY ./internal/config/config.toml /app/config.toml

# Install the CA certificates package
RUN apk --no-cache add ca-certificates

# Expose port 8081
EXPOSE 8081

# Set the environment variables
ENV POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=password \
    POSTGRES_DB=movie_db \
    POSTGRES_HOST=db \
    POSTGRES_PORT=5432

ENV CONFIG_FILE=/app/config.toml

# Start the application
CMD ["/app/main"]
