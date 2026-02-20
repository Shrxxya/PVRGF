package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Criteria struct {
	MinLength      int    `json:"min_length"`
	MinUppercase   int    `json:"min_uppercase"`
	MinLowercase   int    `json:"min_lowercase"`
	MinNumbers     int    `json:"min_numbers"`
	MinSpecial     int    `json:"min_special"`
	AllowedSpecial string `json:"allowed_special"`
}

func LoadCriteria(website string) (*Criteria, error) {
	filePath := "./data/" + strings.ToLower(website) + ".json"

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("website criteria not found")
	}

	var criteria Criteria
	err = json.Unmarshal(file, &criteria)
	if err != nil {
		return nil, err
	}

	return &criteria, nil
}
