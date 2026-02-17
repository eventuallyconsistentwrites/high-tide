package middleware

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin"
)

// RateLimiter is a middleware that uses a Count-Min Sketch to limit requests from IPs.
type RateLimiter struct {
	counter   countmin.BaseCounter
	threshold int
	logger    *slog.Logger
}

// NewRateLimiter creates a new RateLimiter middleware.
func NewRateLimiter(baseCounter *countmin.BaseCounter, threshold int, logger *slog.Logger) *RateLimiter {
	return &RateLimiter{
		counter:   *baseCounter,
		threshold: threshold,
		logger:    logger,
	}
}

// Limit is the middleware handler that wraps another http.Handler.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a production environment, you'd often be behind a reverse proxy.
		// The X-Forwarded-For header is the standard way to identify the originating IP.
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			// If the header is not present, fall back to the remote address.
			var err error
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr // Fallback for addresses without a port
			}
		}

		// Create a logger with context about this specific request
		requestLogger := rl.logger.With("ip", ip, "remote_addr", r.RemoteAddr)
		requestLogger.Info("checking rate limit")

		rl.counter.Update(ip)
		count := rl.counter.PointQuery(ip)

		// rl.counter.DisplayCMS() // This can be very verbose, consider removing from hot path
		if count > rl.threshold {
			requestLogger.Warn("rate limit exceeded", "count", count, "threshold", rl.threshold)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
