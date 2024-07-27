package checker

import (
	"fmt"
	"sync"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

var CheckFunctions []func(string, IPChecker) ([]report.Record, error)
var ProcessWg sync.WaitGroup

type IPChecker interface {
	IsIpReachable(ip string) error
	GetCnameRecords(domain string) (cname string, err error)
	IsCNAMEReachable(cname string) error
	ParseHtml(domain string) (string, error)
	GetDmarcRecords(domain string) ([]string, error)
	ParseSPF(spf string) (map[string]string, bool, string)
	GetTxtRecords(domain string) ([]string, error)
}

type DefaultIPChecker struct{}

func CheckRecords(domainName string, d *report.Data) {
	var wg sync.WaitGroup
	recordChannel := make(chan report.Record)
	ipc := &DefaultIPChecker{}
	ProcessWg.Add(1)

	go func() {
		wg.Wait()
		close(recordChannel)
	}()

	// Simulate record checking
	for _, function := range CheckFunctions {
		wg.Add(1)
		go func(f func(string, IPChecker) ([]report.Record, error)) {
			defer wg.Done()
			// Simulated record checking logic
			record, err := f(domainName, ipc)
			if err != nil {
				fmt.Println("Error checking records:", err)
				return
			}

			for _, r := range record {
				recordChannel <- r
			}
			// wg.Wait()
		}(function)
	}
	go func() {
		defer ProcessWg.Done()
		for record := range recordChannel {
			d.AddRecord(domainName, record.RecordName, record.Status, record.Value, record.Info)
		}
	}()

}

func RegisterCheckFunction(f func(string, IPChecker) ([]report.Record, error)) {
	CheckFunctions = append(CheckFunctions, f)
}

func init() {

	RegisterCheckFunction(Validate_a)
	RegisterCheckFunction(validate_cname)
	RegisterCheckFunction(validate_dkim)
	RegisterCheckFunction(validate_dmarc)
	RegisterCheckFunction(validate_spf)

}
