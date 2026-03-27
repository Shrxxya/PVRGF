package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	baseUrl := "http://localhost:8080/api"
	
	// 1. Register a test user
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())
	password := "TestPass123!"
	
	payload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	
	fmt.Printf("Registering user: %s\n", username)
	resp, err := http.Post(baseUrl+"/register", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to register: %v\n", err)
		return
	}
	resp.Body.Close()
	
	// 2. Login to get UserID
	resp, err = http.Post(baseUrl+"/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to login: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var loginData struct {
		UserId int    `json:"userId"`
		Token  string `json:"token"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &loginData)
	userId := loginData.UserId
	token := loginData.Token
	fmt.Printf("Logged in! UserID: %d, Token: %s...\n", userId, token[:10])

	// 3. Stress test parallel password saving
	passwords := []struct {
		Label string
		Pass  string
	}{
		{"gmail", "Password123!"},
		{"facebook", "Pass!234Word"},
		{"instagram", "Insta!Gram2024"},
		{"linkedin", "Linked!In2024"},
		{"github", "Git!Hub2024"},
	}

	fmt.Println("\nStarting Go Stress Test (5 parallel requests)...")
	start := time.Now()

	var wg sync.WaitGroup
	for _, p := range passwords {
		wg.Add(1)
		go func(label, password string) {
			defer wg.Done()
			
			payload, _ := json.Marshal(map[string]interface{}{
				"label":    label,
				"password": password,
			})

			req, _ := http.NewRequest("POST", baseUrl+"/passwords", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error for %s: %v\n", label, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Error for %s: Status %d - %s\n", label, resp.StatusCode, string(body))
			} else {
				fmt.Printf("Success for %s\n", label)
			}
		}(p.Label, p.Pass)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("\nTotal Time Taken: %v\n", duration)
	fmt.Printf("Sequential Expected: 2.5s\n")
	fmt.Printf("Efficiency Gain: %v\n", 2500*time.Millisecond-duration)
}
