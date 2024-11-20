# Use the official Golang image for the build stage
FROM golang:1.22.4 AS builder

# Set the working directory in the container
WORKDIR /code

# Copy the entire project into the container
COPY . .

# Build the Go binary for your API
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

# Final stage: use a minimal image to run the server
FROM alpine:3.6

# Install necessary dependencies (for certificates)
RUN apk add --no-cache ca-certificates

# Create a user and group to run the app (for better security)
RUN addgroup -g 1001 -S surfe && adduser -u 1001 -S surfe -G surfe

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /bin/server /app/

# Copy environment files
COPY .env.integration /app/.env

# Use a non-root user for security
USER surfe

# Expose the application port
EXPOSE 8081

# Define the default command to run your API
ENTRYPOINT ["/app/server"]
