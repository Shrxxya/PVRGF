package api

import (
	"encoding/json"
	"testing"
	"PVRGF/internal/menu"
)

func TestJSONMarshalUnmarshal(t *testing.T) {
	// 1. Test Struct to JSON (Marshal)
	entry := menu.PasswordEntry{
		ID:        1,
		Label:     "test-site",
		Password:  "Secret123!",
		CreatedAt: "2023-10-27",
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal struct: %v", err)
	}

	// 2. Test JSON to Struct (Unmarshal)
	var decodedEntry menu.PasswordEntry
	err = json.Unmarshal(jsonData, &decodedEntry)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// 3. Verify values
	if decodedEntry.Label != entry.Label {
		t.Errorf("Expected label %s, got %s", entry.Label, decodedEntry.Label)
	}
	if decodedEntry.Password != entry.Password {
		t.Errorf("Expected password to match")
	}
}
