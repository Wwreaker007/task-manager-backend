# 1. Build stage
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /task-manager-backend

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go binaries orders, kitchen, cookhouse
RUN go build -o main .

# 2. Final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /task-manager-backend

# Copy the Go binary from the build stage
COPY --from=builder /task-manager-backend/main .

# Command to run the Go app
CMD ["./main"]