# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages for building (e.g., git for fetching dependencies, ffmpeg for transcoding)
RUN apk add --no-cache git build-base ffmpeg

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build -o /server ./cmd/server/main.go

# Stage 2: Create the final image with the built application
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Install necessary runtime dependencies (e.g., ffmpeg for transcoding if needed during runtime)
RUN apk add --no-cache ffmpeg ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /server .

# Copy static files, if any
COPY web/static ./static

# Set environment variables if needed
# ENV VAR_NAME value

# Expose the port the application will run on
EXPOSE 8080

# Run the binary
CMD ["./server"]