package checker

import (
	"testing"
)

func TestValidateCNAME(t *testing.T) {
	mockChecker := &MockIPChecker{
		CnameRecords: map[string]string{
			"example.com": "cname.example.com",
		},
		ReachableIPs: map[string]bool{
			"cname.example.com": true,
		},
	}

	records, err := validate_cname("example.com", mockChecker)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, record := range records {
		if record.Status == "Fail" && record.Info == "CNAME is not reachable" {
			t.Logf("Failed CNAME: %v", record.Value)
		} else if record.Status == "Pass" {
			t.Logf("Passed CNAME: %v", record.Value)
		}
	}
}
