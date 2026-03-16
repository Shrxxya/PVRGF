package main

import (
	"fmt"
	"net/http"

	"PVRGF/internal/api"
	"PVRGF/internal/db"
)

func main() {
	storage := db.InitDB()
	defer storage.Close()

	handler := api.NewAPIHandler(storage)

	// API Routes
	http.HandleFunc("/api/register", handler.Register)
	http.HandleFunc("/api/login", handler.Login)
	http.HandleFunc("/api/passwords", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPasswords(w, r)
		} else if r.Method == http.MethodPost {
			handler.SavePassword(w, r)
		}
	})
	http.HandleFunc("/api/generate", handler.GeneratePassword)
	http.HandleFunc("/api/criteria", handler.GetCriteria)

	// Serve Static Files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("VAULTSEC Web UI running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
