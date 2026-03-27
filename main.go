package main

import (
	"context"
	"net/http"

	"PVRGF/internal/api"
	"PVRGF/internal/concurrency"
	"PVRGF/internal/db"
	"PVRGF/internal/logger"
	"PVRGF/internal/middleware"
)

var mainLog = logger.New("INFO")

func main() {
	mainLog.Info("Starting VAULTSEC Web Application", nil)
	storage := db.InitDB()
	defer storage.Close()

	// Enhanced Concurrency Worker Pool with Context
	ctx := context.Background()
	concurrency.StartWorkerPool(5, ctx)
	defer concurrency.StopWorkerPool()

	handler := api.NewAPIHandler(storage)

	// API Routes (Unprotected)
	http.HandleFunc("/api/register", handler.Register)
	http.HandleFunc("/api/login", handler.Login)

	// API Routes (Protected by JWTAuth Middleware)
	http.HandleFunc("/api/passwords", middleware.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPasswords(w, r)
		} else if r.Method == http.MethodPost {
			handler.SavePassword(w, r)
		}
	}))
	http.HandleFunc("/api/generate", middleware.JWTAuth(handler.GeneratePassword))
	http.HandleFunc("/api/criteria", middleware.JWTAuth(handler.GetCriteria))
	http.HandleFunc("/api/entropy", middleware.JWTAuth(handler.GetEntropy))

	// Serve Static Files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	mainLog.Info("VAULTSEC Web UI running at http://localhost:8080", nil)

	// Global Logging Middleware
	wrappedMux := middleware.LoggingMiddleware(http.DefaultServeMux)

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		mainLog.Error("Server crashed", map[string]interface{}{"error": err.Error()})
	}
}

