# Start from the latest golang base image
FROM golang:1.25-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

######## Start a new stage from scratch #######
FROM alpine:latest  

WORKDIR /root/

# Install tzdata for timezone support
RUN apk --no-cache add tzdata

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy public directory for static files
COPY --from=builder /app/public ./public

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
