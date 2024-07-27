package checker

import (
	"fmt"
	"net"
	"strings"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

// validate_spf validates the SPF record for a given domain.
func validate_spf(domain string, checker IPChecker) ([]report.Record, error) {
	txtRecords, err := checker.GetTxtRecords(domain)
	if err != nil {
		return []report.Record{}, err
	}

	var records []report.Record
	for _, txt := range txtRecords {

		if !strings.HasPrefix(txt, "v=spf1") {
			continue
		}

		parsedSPF, isValid, reason := checker.ParseSPF(txt)
		status := "Pass"
		if !isValid {
			status = "Fail"
		}

		records = append(records, report.Record{
			RecordName: "SPF",
			Status:     status,
			Value:      parsedSPF,
			Info:       reason,
		})
	}

	if len(records) == 0 {
		return []report.Record{}, nil
	}

	return records, nil
}

// parseSPF parses the SPF record into its components and checks its validity.
func (c *DefaultIPChecker) ParseSPF(spf string) (map[string]string, bool, string) {
	components := strings.Fields(spf)
	spfMap := make(map[string]string)
	reason := ""
	isValid := true
	if strings.Contains(spf, "v=spf") {
		spfMap["v"] = "spf1"
		// isValid = true
	} else {
		reason = "v=spf1 not found"
		return spfMap, false, reason
	}
	if !strings.Contains(spf, "include") {
		isValid = false
		reason = "no authorised domain included"
	}

	for i, comp := range components[1:] {
		key := "param" + fmt.Sprint(i)

		splitbyequal := strings.Split(comp, "=")
		splitbycolon := strings.Split(comp, ":")
		if len(splitbyequal) == 2 {
			spfMap[splitbyequal[0]] = splitbyequal[1]
		} else if len(splitbycolon) == 2 {
			spfMap[splitbycolon[0]] = splitbycolon[1]
		} else {
			spfMap[key] = comp
		}
	}

	return spfMap, isValid, reason
}

func (c *DefaultIPChecker) GetTxtRecords(domain string) ([]string, error) {
	return net.LookupTXT(domain)
}
