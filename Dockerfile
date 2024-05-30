# Use the official Golang image as a build stage
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set environment variables
ARG PORT
ENV PORT=$PORT
ARG DB_DSN
ENV DB_DSN=$DB_DSN
ARG MODE
ENV MODE=$MODE

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main /main

# Expose the port
EXPOSE 5000

# Command to run the executable
CMD ["./main"]