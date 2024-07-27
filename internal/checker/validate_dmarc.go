package checker

import (
	"net"
	"strings"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

func validate_dmarc(domain string, checker IPChecker) ([]report.Record, error) {
	var records []report.Record
	valueMap := make(map[string]string)
	var reason string = ""
	status := "Pass"
	resp, err := checker.GetDmarcRecords(domain)
	if err != nil {
		records = append(records, report.Record{
			RecordName: "DMARC",
			Status:     "Fail",
			Value:      "",
			Info:       "DMARC not found",
		})
		return records, nil
	}
	for _, txt := range resp {
		if !strings.Contains(txt, "v=DMARC") {
			status = "warn"
			reason = "v=DMARC1 not found"
		}
		if !strings.Contains(txt, "; p=") {
			status = "warn"
			reason = reason + "; policy not found"
		}
		if !strings.Contains(txt, "; rua=") {
			status = "warn"
			reason = reason + "; reporting URI not found"
		}
		if !strings.Contains(txt, "; ruf=") {
			status = "warn"
			reason = reason + "; reporting address not found"
		}

		// split the txt record by spacke
		txtSplit := strings.Fields(txt)
		for _, t := range txtSplit {
			if strings.Contains(t, "=") {
				split := strings.Split(t, "=")
				valueMap[split[0]] = split[1]
			}
		}
		records = append(records, report.Record{
			RecordName: "DMARC",
			Status:     status,
			Value:      valueMap,
			Info:       reason,
		})

	}
	return records, nil
}

func (c *DefaultIPChecker) GetDmarcRecords(domain string) (dmarcs []string, err error) {
	resp, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return []string{}, err
	}
	return resp, nil
}
