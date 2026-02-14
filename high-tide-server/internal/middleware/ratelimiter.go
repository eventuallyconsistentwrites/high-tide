package middleware

import (
	"log"
	"net"
	"net/http"

	"github.com/eventuallyconsistentwrites/high-tide-server/countmin"
)

// RateLimiter is a middleware that uses a Count-Min Sketch to limit requests from IPs.
type RateLimiter struct {
	cms       *countmin.CountMinSketch
	threshold int
}

// NewRateLimiter creates a new RateLimiter middleware.
func NewRateLimiter(cms *countmin.CountMinSketch, threshold int) *RateLimiter {
	return &RateLimiter{
		cms:       cms,
		threshold: threshold,
	}
}

// Limit is the middleware handler that wraps another http.Handler.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Checking Rate Limit for Remote Address %s\n", r.RemoteAddr)
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // Fallback for addresses without a port
		}
		log.Printf("Checking Rate Limit for IP %s\n", ip)
		rl.cms.Update(ip)
		count := rl.cms.PointQuery(ip)
		log.Printf("IP: %s, count: %d\n", ip, count)
		rl.cms.DisplayCMS()
		if count > rl.threshold {
			log.Printf("Rate limit exceeded for IP %s with count %d\n", ip, count)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
