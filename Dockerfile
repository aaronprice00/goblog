FROM golang:alpine AS builder

LABEL maintainer="Aaron Price <aaronprice00@gmail.com>"

# Install Git
RUN apk update && apk add --no-cache git

# Set CWD in container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source files
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage (from scratch)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy built binary && .env from previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose the port
EXPOSE 8080

# Run the executable
CMD ["./main"]