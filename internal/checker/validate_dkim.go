// package checker

// import (
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/anaskhan96/soup"
// 	"github.com/mukul-kr/dns-verifier/pkg/report"
// )

// func validate_dkim(domain string) ([]report.Record, error) {
// 	var records []report.Record
// 	var reason string = ""
// 	valueMap := make(map[string]string)
// 	status := "Pass"
// 	dkim, err := parseHtml(domain)
// 	if err != nil || dkim == "" {
// 		if err.Error() == "Dkim not found" {
// 			records = append(records, report.Record{
// 				RecordName: "DKIM",
// 				Status:     "Fail",
// 				Value:      "",
// 				Info:       "DKIM record not found",
// 			})
// 		} else {
// 			return []report.Record{}, err
// 		}
// 	}
// 	if !strings.Contains(dkim, "DKIM1") {
// 		status = "warn"
// 		reason = "DKIM1 not found"
// 	}
// 	if !strings.Contains(dkim, "; k=") {
// 		status = "warn"
// 		reason = reason + "; key not found"
// 	}
// 	if !strings.Contains(dkim, "; p=") {
// 		status = "warn"
// 		reason = reason + "; public key not found"
// 	}
// 	txtSplit := strings.Fields(dkim)
// 	for _, t := range txtSplit {
// 		if strings.Contains(t, "=") {
// 			split := strings.Split(t, "=")
// 			valueMap[split[0]] = split[1]
// 		}
// 	}
// 	if reason != "" {
// 		records = append(records, report.Record{
// 			RecordName: "DKIM",
// 			Status:     status,
// 			Value:      valueMap,
// 			Info:       reason,
// 		})
// 	} else {
// 		records = append(records, report.Record{
// 			RecordName: "DKIM",
// 			Status:     status,
// 			Value:      dkim,
// 		})
// 	}
// 	return records, nil

// }

// func parseHtml(domain string) (string, error) {
// 	resp, err := soup.Get(fmt.Sprintf("https://registry.prove.email/?domain=%s", domain))
// 	if err != nil {
// 		fmt.Println("Error fetching the page")
// 		return "", errors.New("Error fetching the page")
// 	}
// 	doc := soup.HTMLParse(resp)

// 	pre := doc.Find("pre")
// 	if pre.Error != nil {
// 		fmt.Println("Error fetching the pre tag", pre.Error)
// 		return "", errors.New("Dkim not found")
// 	}

// 	return pre.Text(), nil
// }

package checker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/mukul-kr/dns-verifier/pkg/report"
)

var errDkim = errors.New("dkim not found")
var errFetch = errors.New("error fetching the page")

func validate_dkim(domain string, checker IPChecker) ([]report.Record, error) {
	var records []report.Record
	var reason string = ""
	valueMap := make(map[string]string)
	status := "Pass"

	dkim, err := checker.ParseHtml(domain)

	if err != nil || dkim == "" {
		if err != nil && errors.Is(err, errDkim) {
			reason = "DKIM record not found"
			records = append(records, report.Record{
				RecordName: "DKIM",
				Status:     "Fail",
				Value:      "",
				Info:       reason,
			})
			return records, nil
		} else {
			return []report.Record{}, err
		}
	}
	if !strings.Contains(dkim, "DKIM1") {
		status = "warn"
		reason = "DKIM1 not found"
	}
	if !strings.Contains(dkim, "; k=") {
		status = "warn"
		reason = reason + "; key not found"
	}
	if !strings.Contains(dkim, "; p=") {
		status = "warn"
		reason = reason + "; public key not found"
	}
	txtSplit := strings.Fields(dkim)
	for _, t := range txtSplit {
		if strings.Contains(t, "=") {
			split := strings.Split(t, "=")
			valueMap[split[0]] = split[1]
		}
	}
	records = append(records, report.Record{
		RecordName: "DKIM",
		Status:     status,
		Value:      valueMap,
		Info:       reason,
	})
	return records, nil
}

func (c *DefaultIPChecker) ParseHtml(domain string) (string, error) {
	resp, err := soup.Get(fmt.Sprintf("https://registry.prove.email/?domain=%s", domain))
	if err != nil {
		// fmt.Println("Error fetching the page")
		return "", errFetch
	}
	doc := soup.HTMLParse(resp)

	pre := doc.Find("pre")
	if pre.Error != nil {
		// fmt.Println("Error fetching the pre tag", pre.Error)
		return "", errDkim
	}
	return pre.Text(), nil
}
