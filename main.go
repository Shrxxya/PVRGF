package main

import (
	"fmt"
	"net/http"

	"PVRGF/internal/api"
	"PVRGF/internal/concurrency"
	"PVRGF/internal/db"
)

func main() {
	storage := db.InitDB()
	defer storage.Close()

	// Initialize Concurrency Worker Pool (Educational: Goroutines)
	concurrency.StartWorkerPool(5)

	handler := api.NewAPIHandler(storage)

	// API Routes (Protected by Middleware)
	http.HandleFunc("/api/register", api.RateLimitMiddleware(handler.Register))
	http.HandleFunc("/api/login", api.RateLimitMiddleware(handler.Login))

	http.HandleFunc("/api/passwords", api.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPasswords(w, r)
		} else if r.Method == http.MethodPost {
			handler.SavePassword(w, r)
		}
	}))

	http.HandleFunc("/api/generate", api.JWTMiddleware(handler.GeneratePassword))
	http.HandleFunc("/api/criteria", api.JWTMiddleware(handler.GetCriteria))

	// Serve Static Files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("VAULTSEC Web UI running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
