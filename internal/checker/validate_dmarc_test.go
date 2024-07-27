package checker

import (
	"testing"
)

func TestValidateDMARC(t *testing.T) {
	mockChecker := &MockIPChecker{
		DmarcRecords: map[string][]string{
			"example.com": {
				"v=DMARC1; p=none; rua=mailto:dmarc-reports@example.com; ruf=mailto:dmarc-failures@example.com",
			},
		},
	}

	records, err := validate_dmarc("example.com", mockChecker)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, record := range records {
		t.Logf("Record: %+v", record)
	}
}
