package checker

import (
	"testing"
)

func TestValidateSPF(t *testing.T) {
	mockChecker := &MockIPChecker{
		TxtRecords: map[string][]string{
			"example.com": {
				"v=spf1 include:_spf.google.com ~all",
			},
		},
		ParsedSPF: map[string]map[string]string{
			"v=spf1 include:_spf.google.com ~all": {
				"v":       "spf1",
				"include": "_spf.google.com",
				"param0":  "~all",
			},
		},
	}

	records, err := validate_spf("example.com", mockChecker)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, record := range records {
		t.Logf("Record: %+v", record)
	}
}
