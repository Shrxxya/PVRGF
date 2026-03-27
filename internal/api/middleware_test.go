package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"PVRGF/internal/auth"
)

func TestJWTMiddleware(t *testing.T) {
	// 1. Generate a valid token
	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 2. Setup handler and middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey).(int)
		if userID != 1 {
			t.Errorf("Expected userID 1 in context, got %v", userID)
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := JWTMiddleware(nextHandler)

	// 3. Test valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	// 4. Test missing token
	req = httptest.NewRequest("GET", "/test", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized for missing token, got %v", rr.Code)
	}

	// 5. Test invalid token
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized for invalid token, got %v", rr.Code)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := RateLimitMiddleware(nextHandler)

	// Test 5 requests (should be allowed due to burst)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/login", nil)
		// Use a unique RemoteAddr to avoid interference if running tests multiple times
		req.RemoteAddr = "1.2.3.4:1234"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status OK, got %v", i+1, rr.Code)
		}
	}

	// 6th request should be blocked
	req := httptest.NewRequest("POST", "/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status TooManyRequests for 6th request, got %v", rr.Code)
	}
}
