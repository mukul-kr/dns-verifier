package checker

import (
	"testing"
)

func TestValidateA(t *testing.T) {
	mockChecker := &MockIPChecker{
		ReachableIPs: map[string]bool{
			"192.168.1.1": true,
			"192.168.1.2": false,
		},
	}

	// You can also mock net.LookupIP for full isolation if needed
	records, err := Validate_a("example.com", mockChecker)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, record := range records {
		if record.Status == "Fail" && record.Info == "IP address is not reachable" {
			t.Logf("Failed IP: %v", record.Value)
		} else if record.Status == "Pass" {
			t.Logf("Passed IP: %v", record.Value)
		}
	}
}
