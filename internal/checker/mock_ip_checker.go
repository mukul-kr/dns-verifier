// mock_ip_checker.go
package checker

import "errors"

// MockIPChecker is a mock implementation of the IPChecker interface for testing purposes.
type MockIPChecker struct {
	ReachableIPs map[string]bool
	CnameRecords map[string]string
	DmarcRecords map[string][]string
	TxtRecords   map[string][]string
	ParsedSPF    map[string]map[string]string
	HtmlRecords  map[string]string
}

func (m *MockIPChecker) IsIpReachable(ip string) error {
	if reachable, exists := m.ReachableIPs[ip]; exists && reachable {
		return nil
	}
	return errors.New("IP not reachable")
}

func (m *MockIPChecker) GetCnameRecords(domain string) (string, error) {
	if cname, exists := m.CnameRecords[domain]; exists {
		return cname, nil
	}
	return "", errors.New("CNAME record not found")
}

func (m *MockIPChecker) IsCNAMEReachable(cname string) error {
	if reachable, exists := m.ReachableIPs[cname]; exists && reachable {
		return nil
	}
	return errors.New("CNAME not reachable")
}

func (m *MockIPChecker) ParseHtml(domain string) (string, error) {
	if html, exists := m.HtmlRecords[domain]; exists {
		return html, nil
	}
	return "", errors.New("HTML not found")
}

func (m *MockIPChecker) GetDmarcRecords(domain string) ([]string, error) {
	if dmarc, exists := m.DmarcRecords[domain]; exists {
		return dmarc, nil
	}
	return nil, errors.New("DMARC records not found")
}

func (m *MockIPChecker) ParseSPF(spf string) (map[string]string, bool, string) {
	if parsed, exists := m.ParsedSPF[spf]; exists {
		return parsed, true, ""
	}
	return nil, false, "SPF parsing failed"
}

func (m *MockIPChecker) GetTxtRecords(domain string) ([]string, error) {
	if txt, exists := m.TxtRecords[domain]; exists {
		return txt, nil
	}
	return nil, errors.New("TXT records not found")
}
