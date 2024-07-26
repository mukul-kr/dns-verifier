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
				fmt.Println(fmt.Sprintf("Checking %s record for %s", r.RecordName, domainName), "Status", r.Status, "Value", r.Value)
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
	// l, _ := logger.NewLogger()

	RegisterCheckFunction(validate_a)
}

// func validate_a(domainName string) ([]report.Record, error) {
// 	return []report.Record{
// 		{RecordName: "A", Status: "ok", Value: "1.1.1.1"},
// 	}, nil
// }
