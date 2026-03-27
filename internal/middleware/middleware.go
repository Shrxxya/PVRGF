package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"PVRGF/internal/logger"
)

var (
	// In production, this would come from an ENV var
	jwtSecretKey = []byte("my_super_secret_vaultsec_key_2026")
	appLogger    = logger.New("INFO")
)

// JWTAuth middleware verifies the JWT token present in the Authorization header
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// String manipulation: extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization Header Format", http.StatusUnauthorized)
			return
		}
		
		tokenString := parts[1]
		
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired JWT Token", http.StatusUnauthorized)
			return
		}
		
		// Extract userID (float64 due to JSON parsing, convert to int)
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid Token Payload", http.StatusUnauthorized)
			return
		}
		
		// Add to context
		ctx := context.WithValue(r.Context(), "userID", int(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GenerateJWT creates a new JWT for a successful login
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// LoggingMiddleware logs all incoming requests utilizing the new custom logger
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		appLogger.Info("API Request", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
			"client":   r.RemoteAddr,
		})
	})
}
