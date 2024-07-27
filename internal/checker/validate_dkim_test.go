// validate_dkim_test.go
package checker

import (
	"testing"
)

func TestValidateDKIM(t *testing.T) {
	mockChecker := &MockIPChecker{
		HtmlRecords: map[string]string{
			"example.com": "v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA",
		},
	}

	records, err := validate_dkim("example.com", mockChecker)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, record := range records {
		t.Logf("Record: %+v", record)
	}
}
