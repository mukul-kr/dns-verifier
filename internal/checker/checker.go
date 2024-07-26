package checker

import (
	"fmt"
	"sync"

	"github.com/mukul-kr/dns-verifier/pkg/report"
)

var checkFunctions []func(string) ([]report.Record, error)
var ProcessWg sync.WaitGroup

// var log *zap.SugaredLogger

func CheckRecords(domainName string, d *report.Data) {
	var wg sync.WaitGroup
	recordChannel := make(chan report.Record)
	ProcessWg.Add(1)

	go func() {
		wg.Wait()
		close(recordChannel)
	}()

	// Simulate record checking
	for _, function := range checkFunctions {
		wg.Add(1)
		go func(f func(string) ([]report.Record, error)) {
			defer wg.Done()
			// Simulated record checking logic
			record, err := f(domainName)
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
			d.AddRecord(domainName, record.RecordName, record.Status, record.Value)
		}
	}()

}

func RegisterCheckFunction(f func(string) ([]report.Record, error)) {
	checkFunctions = append(checkFunctions, f)
}

func init() {

	RegisterCheckFunction(validate_a)
	RegisterCheckFunction(validate_cname)
	RegisterCheckFunction(validate_dkim)
	RegisterCheckFunction(validate_dmarc)
	RegisterCheckFunction(validate_spf)

}
