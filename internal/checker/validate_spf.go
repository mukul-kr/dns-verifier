package checker

import (
	"net"
	"strings"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

// validate_spf validates the SPF record for a given domain.
func validate_spf(domain string) ([]report.Record, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return []report.Record{}, err
	}

	var records []report.Record
	for _, txt := range txtRecords {

		parsedSPF, isValid, reason := parseSPF(txt, domain)
		status := "Pass"
		if !isValid {
			status = "Fail"
		}
		if reason != "" {
			records = append(records, report.Record{
				RecordName: "SPF",
				Status:     status,
				Value:      parsedSPF,
			})
		} else {
			records = append(records, report.Record{
				RecordName: "SPF",
				Status:     status,
				Value:      parsedSPF,
				Info:       reason,
			})
		}

	}

	if len(records) == 0 {
		return []report.Record{}, nil
	}

	return records, nil
}

// parseSPF parses the SPF record into its components and checks its validity.
func parseSPF(spf, domain string) (map[string]string, bool, string) {
	components := strings.Fields(spf)
	spfMap := make(map[string]string)
	reason := ""
	isValid := false
	if components[0] == "v=spf1" {
		spfMap["v"] = "spf1"
		isValid = true
	} else {
		reason = "v=spf1 not found"
		return spfMap, false, reason
	}

	for i, comp := range components[1:] {
		key := "param" + string(i)
		if !strings.Contains(comp, "include") {
			isValid = false
			reason = "no authorised domain included"
		}
		spfMap[key] = comp
	}

	return spfMap, isValid, reason
}
