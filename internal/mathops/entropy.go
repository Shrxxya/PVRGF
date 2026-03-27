package mathops

import (
	"math"
	"strings"
)

// CalculateEntropy performs mathematical operations to calculate passwords True Entropy.
// Formula: E = L * log2(R)
// Where L is length, R is size of the character pool
func CalculateEntropy(password string) float64 {
	L := float64(len(password))
	R := 0.0
	
	hasLower := false
	hasUpper := false
	hasDigits := false
	hasSpecial := false
	
	// String manipulation to check character presence
	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigits = true
		default:
			if strings.ContainsRune("!@#$%^&*()-_+=~[]{}\\|:;\"'<>,.?/", char) {
				hasSpecial = true
			}
		}
	}
	
	if hasLower { R += 26 }
	if hasUpper { R += 26 }
	if hasDigits { R += 10 }
	if hasSpecial { R += 32 }
	
	if R == 0 {
		return 0
	}
	
	// Core Mathematical operations
	// E = L * log2(R)
	entropy := L * math.Log2(R)
	
	// Round to 2 decimal places using math.Round
	return math.Round(entropy*100) / 100
}

// EvaluateStrength generates a user-friendly mathematical strength evaluation based on entropy
func EvaluateStrength(entropy float64) string {
	if entropy < 28 {
		return "Very Weak"
	}
	if entropy < 36 {
		return "Weak"
	}
	if entropy < 60 {
		return "Reasonable"
	}
	if entropy < 128 {
		return "Strong"
	}
	return "Very Strong"
}
