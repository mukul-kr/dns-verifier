package domain

import (
	"net"
	"net/url"
)

func parseDomain(domain string) (string, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return "", err
	}
	host, _, _ := net.SplitHostPort(u.Host)

	return host, nil
}
