# --- Build Stage ---
# Use the official Go image as a builder.
FROM golang:1.23-alpine AS builder

# Install build tools and SQLite dev libraries required for CGO.
RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies and leverage Docker cache.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Tidy the modules to ensure go.sum is correct before building.
RUN go mod tidy

# Build the binary with CGO enabled to support go-sqlite3.
RUN go build -o /server ./cmd/api/

# --- Final Stage ---
# Use a minimal, non-root base image for the final container.
FROM alpine:latest
WORKDIR /app
# The go-sqlite3 driver requires runtime libraries, which are present in the base alpine image.
COPY --from=builder /server .
EXPOSE 8080
CMD ["./server"]