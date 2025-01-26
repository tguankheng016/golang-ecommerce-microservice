# Use the official Golang image to create a build artifact.
FROM golang:1.23.4 AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux

# Set the Current Working Directory inside the container to pkg
WORKDIR /app/pkg

# Copy pkg go mod and sum files
COPY internal/pkg/go.mod internal/pkg/go.sum ./

# Download dependencies of pkg
RUN go mod download

# Set the Current Working Directory inside the container to product_service
WORKDIR /app/services/product_service

# Copy product_service go mod and sum files
COPY internal/services/product_service/go.mod internal/services/product_service/go.sum ./

# Download dependencies of product_service
RUN go mod download

# Set back to app directory
WORKDIR /app

# Copy the source from pkg and product_service
COPY ./internal .

# Set to product_service directory
WORKDIR /app/services/product_service

# Build the Go app
RUN go build -o main ./cmd/app

# Start a new stage from scratch
FROM alpine:latest  

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/services/product_service/main .

# Copy the config.json files
COPY --from=builder /app/services/product_service/config /config

# Copy the migrations
COPY --from=builder /app/services/product_service/internal/data/migrations /internal/data/migrations

# Copy the go mod for locate config.json
COPY --from=builder /app/services/product_service/go.mod .

# Expose port 8000 to the outside world
EXPOSE 8000
EXPOSE 6600

# Command to run the executable
CMD ["./main"]