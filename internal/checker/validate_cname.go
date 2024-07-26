package checker

import (
	"net"
	"time"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

func validate_cname(domain string) ([]report.Record, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return []report.Record{}, err
	}

	var records []report.Record
	if cname == domain {
		// If the CNAME record is the same as the domain, it's not a valid CNAME record.
		records = append(records, report.Record{
			RecordName: "CNAME",
			Status:     "Fail",
			Value:      cname,
		})
	} else {
		// Check if the CNAME is reachable
		err := isReachable(cname)
		status := "Pass"
		if err != nil {
			status = "Fail"
		}
		records = append(records, report.Record{
			RecordName: "CNAME",
			Status:     status,
			Value:      cname,
		})
	}

	return records, nil
}

// isReachable checks if the CNAME record is reachable.
func isReachable(cname string) error {
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
