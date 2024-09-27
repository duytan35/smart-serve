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
ENV PORT=$PORT
ENV DB_DSN=$DB_DSN
ENV MODE=$MODE
ENV JWT_SECRET=$JWT_SECRET
ENV S3_ACCESS_KEY=$S3_ACCESS_KEY
ENV S3_SECRET_KEY=$S3_SECRET_KEY
ENV S3_BUCKET_NAME=$S3_BUCKET_NAME
ENV S3_REGION=$S3_REGION
ENV S3_URL=$S3_URL
ENV ADMIN_MAIL=$ADMIN_MAIL

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main /main

# Expose the port
EXPOSE 5000

# Command to run the executable
CMD ["./main"]