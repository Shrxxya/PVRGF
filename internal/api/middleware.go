package api

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"PVRGF/internal/auth"
	"golang.org/x/time/rate"
)

type contextKey string

const UserIDKey contextKey = "userID"

// JWTMiddleware validates the bearer token and adds userID to context
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// getLimiter returns a rate limiter for a given key (IP address)
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		// Allow 5 requests per minute, with a burst of 5
		limiter = rate.NewLimiter(rate.Limit(5.0/60.0), 5)
		limiters[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware limits requests per IP address
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}
