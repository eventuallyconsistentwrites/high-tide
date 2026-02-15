# --- Build Stage ---
# Use the official Go image as a builder.
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
# Build a statically-linked binary to ensure it runs in a minimal image.
RUN CGO_ENABLED=0 go build -o /tester main.go

# --- Final Stage ---
# Use a minimal, non-root base image for the final container.
FROM alpine:latest
COPY --from=builder /tester /tester
CMD ["/tester"]