# Dockerfile for Go Backend

# Stage 1: Build the Go app
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Stage 2: Run the Go app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# Install bash for the sleep command
RUN apk add --no-cache bash

# Set environment variables for PostgreSQL connection
ENV DB_HOST=${DB_HOST}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV DB_PORT=${DB_PORT}

# Expose port
EXPOSE 9000

# Command to run the executable with a delay
CMD ["bash", "-c", "sleep 10 && ./main"]
