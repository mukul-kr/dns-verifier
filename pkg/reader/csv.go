package reader

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type CSVHandler struct{}

func (h CSVHandler) Handle(content string) func() ([]string, error) {
	return func() ([]string, error) {
		// Implement CSV handling logic here

		reader := csv.NewReader(strings.NewReader(content))

		records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV data: %w", err)
		}

		// Flatten the records into a single list of strings
		var result []string
		for _, record := range records {
			result = append(result, record...)
		}
		return result, nil

	}
}
