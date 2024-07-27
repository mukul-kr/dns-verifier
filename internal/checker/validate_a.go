package checker

import (
	"net"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

// IPChecker defines an interface for checking if an IP is reachable.

// DefaultIPChecker is a default implementation of IPChecker using net package.

// IsIpReachable checks if the IP is reachable by attempting to establish a TCP connection.
func (c *DefaultIPChecker) IsIpReachable(ip string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(ip, "80"))
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// validate_a validates the A records for the given domain.
func Validate_a(domain string, checker IPChecker) ([]report.Record, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return []report.Record{}, err
	}
	var records []report.Record
	for _, ip := range ips {
		if net.ParseIP(ip.String()) == nil {
			records = append(records, report.Record{
				RecordName: "A",
				Status:     "Fail",
				Value:      ip.String(),
				Info:       "Invalid IP address",
			})
			return records, nil
		} else {
			err := checker.IsIpReachable(ip.String())
			if err != nil {
				records = append(records, report.Record{
					RecordName: "A",
					Status:     "Fail",
					Value:      ip.String(),
					Info:       "IP address is not reachable",
				})
			} else {
				records = append(records, report.Record{
					RecordName: "A",
					Status:     "Pass",
					Value:      ip.String(),
					Info:       "",
				})
			}
		}
	}
	return records, nil
}
