package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"PVRGF/internal/auth"
	"PVRGF/internal/db"
	"PVRGF/internal/menu"
	"PVRGF/internal/rules"
)

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
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"message": "Login successful",
		"userId":  userID,
	})
	w.Write(resp)
}

func (h *APIHandler) GetPasswords(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

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

	body, _ := io.ReadAll(r.Body)
	var data struct {
		UserID   int    `json:"userId"`
		Label    string `json:"label"`
		Password string `json:"password"`
	}
	json.Unmarshal(body, &data)

	// Server-side validation
	criteria, err := rules.LoadCriteria(data.Label)
	if err == nil {
		if !menu.ValidatePassword(data.Password, criteria) {
			http.Error(w, "Password does not meet criteria", http.StatusBadRequest)
			return
		}
	}

	if err := menu.SavePasswordEntry(h.Store.GetDB(), data.UserID, data.Label, data.Password); err != nil {
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
