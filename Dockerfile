# Use the official Golang image for building
FROM golang:1.27-alpine

# Set working directory
WORKDIR /usr/src/app

# Copy Go modules and dependencies
COPY go.mod go.sum ./
COPY cmd ./cmd/
COPY internal ./internal/

RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]
