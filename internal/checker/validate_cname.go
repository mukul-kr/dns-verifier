package checker

import (
	"net"
	"time"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

func validate_cname(domain string, checker IPChecker) ([]report.Record, error) {
	cname, err := checker.GetCnameRecords(domain)
	var reason string = ""
	if err != nil {
		return []report.Record{}, err
	}

	var records []report.Record

	// Check if the CNAME is reachable
	err = checker.IsCNAMEReachable(cname)
	status := "Pass"
	if err != nil {
		status = "Fail"
		reason = "CNAME is not reachable"
	}
	records = append(records, report.Record{
		RecordName: "CNAME",
		Status:     status,
		Value:      cname,
		Info:       reason,
	})

	return records, nil
}

func (c *DefaultIPChecker) GetCnameRecords(domain string) (cname string, err error) {
	return net.LookupCNAME(domain)
}

// isReachable checks if the CNAME record is reachable.
func (c *DefaultIPChecker) IsCNAMEReachable(cname string) error {
	// Resolve the IP address of the CNAME record
	ips, err := net.LookupIP(cname)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		// Attempt to establish a TCP connection to the IP address on port 80 (HTTP)
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip.String(), "80"), 5*time.Second)
		if err != nil {
			return err
		}
		conn.Close()
	}

	return nil
}
