# --- STAGE 1: Build the Go binary ---
FROM golang:1.25-alpine3.21 AS builder

WORKDIR /app

# Install Git (needed for some Go modules)
RUN apk add --no-cache git

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
ENV GOPROXY=direct
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary with CGO disabled for pristine static linking on minimal Alpine runtimes
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- STAGE 2: Create a minimal image to run ---
FROM alpine:3.21

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Expose the port your Go app uses
EXPOSE 8080

# Run the app
CMD ["./server"]
