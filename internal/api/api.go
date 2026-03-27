package api

import (
	"encoding/json"
	"io"
	"net/http"

	"PVRGF/internal/auth"
	"PVRGF/internal/concurrency"
	"PVRGF/internal/db"
	"PVRGF/internal/logger"
	"PVRGF/internal/mathops"
	"PVRGF/internal/menu"
	"PVRGF/internal/middleware"
	"PVRGF/internal/rules"
)

var appLog = logger.New("INFO")


// APIHandler demonstrates pointer receivers (Educational: Pointers)
type APIHandler struct {
	Store db.VaultStorage
}

func NewAPIHandler(store db.VaultStorage) *APIHandler {
	return &APIHandler{Store: store}
}

func (h *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Educational: JSON Unmarshal
	body, _ := io.ReadAll(r.Body)
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := auth.RegisterUser(h.Store.GetDB(), data.Username, data.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Educational: JSON Marshal
	resp, _ := json.Marshal(map[string]string{"message": "Registration successful"})
	w.Write(resp)
}

func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.Unmarshal(body, &data)

	userID, err := auth.LoginUser(h.Store.GetDB(), data.Username, data.Password)
	if err != nil {
		appLog.Warn("Failed login attempt", map[string]interface{}{"username": data.Username})
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateJWT(userID)
	if err != nil {
		appLog.Error("Token generation failed", map[string]interface{}{"error": err.Error()})
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	
	appLog.Info("User logged in successfully", map[string]interface{}{"userID": userID})

	resp, _ := json.Marshal(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
	w.Write(resp)
}

func (h *APIHandler) GetPasswords(w http.ResponseWriter, r *http.Request) {
	// Secure: Extract userID from JWT Context context, not from URL parameters
	userIDObj := r.Context().Value("userID")
	if userIDObj == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDObj.(int)


	passwords, err := menu.GetPasswords(h.Store.GetDB(), userID)
	if err != nil {
		http.Error(w, "Error retrieving passwords", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(passwords)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *APIHandler) SavePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Secure: Get userID from Context, ignore whatever is in the body
	userIDObj := r.Context().Value("userID")
	if userIDObj == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDObj.(int)

	body, _ := io.ReadAll(r.Body)
	var data struct {
		Label    string `json:"label"`
		Password string `json:"password"`
	}
	json.Unmarshal(body, &data)

	appLog.Info("Saving password attempt", map[string]interface{}{
		"userID": userID,
		"label":  data.Label,
		"password": data.Password, // This will be safely masked!
	})



	// Server-side validation using Goroutine Worker Pool
	if !concurrency.SubmitTask(data.Label, data.Password) {
		appLog.Warn("Password criteria validation failed", map[string]interface{}{"label": data.Label})
		http.Error(w, "Password does not meet criteria (validated via goroutine)", http.StatusBadRequest)
		return
	}

	if err := menu.SavePasswordEntry(h.Store.GetDB(), userID, data.Label, data.Password); err != nil {
		appLog.Error("Failed to save password", map[string]interface{}{"error": err.Error()})
		http.Error(w, "Error saving password", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(map[string]string{"message": "Password saved"})
	w.Write(resp)
}

func (h *APIHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	label := r.URL.Query().Get("label")
	if label == "" {
		label = "gmail"
	}

	pwd, err := menu.GeneratePassword(label)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(map[string]string{"password": pwd})
	w.Write(resp)
}


func (h *APIHandler) GetCriteria(w http.ResponseWriter, r *http.Request) {
	label := r.URL.Query().Get("label")
	if label == "" {
		http.Error(w, "Label is required", http.StatusBadRequest)
		return
	}

	criteria, err := rules.LoadCriteria(label)
	if err != nil {
		http.Error(w, "Criteria not found", http.StatusNotFound)
		return
	}

	resp, _ := json.Marshal(criteria)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *APIHandler) GetEntropy(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	if password == "" {
		http.Error(w, "Password query parameter is required", http.StatusBadRequest)
		return
	}
	
	entropy := mathops.CalculateEntropy(password)
	strength := mathops.EvaluateStrength(entropy)
	
	resp, _ := json.Marshal(map[string]interface{}{
		"entropy":  entropy,
		"strength": strength,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
