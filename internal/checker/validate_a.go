package checker

import (
	"fmt"
	"net"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

func validate_a(domain string) ([]report.Record, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return []report.Record{}, err
	}
	fmt.Println(ips)
	var records []report.Record
	for _, ip := range ips {
		if net.ParseIP(ip.String()) == nil {
			records = append(records, report.Record{
				RecordName: "A",
				Status:     "Fail",
				Value:      ip.String(),
			})
		} else {
			records = append(records, report.Record{
				RecordName: "A",
				Status:     "Pass",
				Value:      ip.String(),
			})
		}
	}
	return records, nil
}
