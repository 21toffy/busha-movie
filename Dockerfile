# Start from the official golang image
FROM golang:1.17-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and cache go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

RUN go build -o app ./cmd/api/main.go

# Start from a new image
FROM alpine:latest

# Install required system libraries
RUN apk --no-cache add ca-certificates

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary file from the build stage
COPY --from=build /app/app .

# Set environment variables for Postgres
ENV POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=password \
    POSTGRES_DB=movies_db \
    POSTGRES_HOST=db \
    POSTGRES_PORT=5432

# Expose port 8080 for the application
EXPOSE 8080

# Start the application
CMD ["./app"]
