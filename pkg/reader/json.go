package reader

import (
	"encoding/json"
	"fmt"
)

type JSONHandler struct{}

func (h JSONHandler) Handle(content string) func() ([]string, error) {
	return func() ([]string, error) {
		// Implement JSON handling logic here
		var data []struct {
			URL string `json:"url"`
		}

		if err := json.Unmarshal([]byte(content), &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON data: %w", err)
		}

		var result []string
		for _, item := range data {
			result = append(result, item.URL)
		}

		return result, nil
	}
}
