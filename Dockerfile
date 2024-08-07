# Use an official Go image as the base image for the build stage
FROM golang:1.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use an official Ubuntu image for the final stage
FROM ubuntu:22.04

# Install MySQL client
RUN apt-get update && apt-get install -y mysql-client && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application will run on
EXPOSE 8001

# Set environment variables for MySQL connection, placeholder will be overridden by docker-compose.yml
ENV GORM_DNS="root:123123@tcp(db:3306)/ordermanager?charset=utf8mb4&parseTime=True&loc=Local"

# Run the Go application
CMD ["./main"]
