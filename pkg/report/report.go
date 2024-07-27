package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Record represents a record with a name, status, and value
type Record struct {
	RecordName string      `json:"recordname" yaml:"recordname"`
	Status     string      `json:"status" yaml:"status"`
	Value      interface{} `json:"value" yaml:"value"`
	Info       string      `json:"info" yaml:"info"`
}

// Domain represents a domain with a name and associated records
type Domain struct {
	DomainName string   `json:"domainname" yaml:"domainname"`
	Records    []Record `json:"records" yaml:"records"`
}

// Data represents the overall structure containing multiple domains
type Data struct {
	Domains []Domain `json:"domains" yaml:"domains"`
	mu      sync.Mutex
}

func (d *Data) AddDomain(domainName string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, domain := range d.Domains {
		if domain.DomainName == domainName {
			return
		}
	}
	d.Domains = append(d.Domains, Domain{
		DomainName: domainName,
		Records:    []Record{},
	})
}

func (d *Data) AddRecord(domainName, recordName, status string, value interface{}, info string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var domain *Domain
	for i := range d.Domains {
		if d.Domains[i].DomainName == domainName {
			domain = &d.Domains[i]
			break
		}
	}
	if domain == nil {
		return
	}
	if len(info) != 0 {
		domain.Records = append(domain.Records, Record{
			RecordName: recordName,
			Status:     status,
			Value:      value,
			Info:       info,
		})
	} else {
		domain.Records = append(domain.Records, Record{
			RecordName: recordName,
			Status:     status,
			Value:      value,
		})
	}
}

// ToJSON converts the Data struct to a JSON string
func (d *Data) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToYAML converts the Data struct to a YAML string
func (d *Data) ToYAML() (string, error) {
	yamlData, err := yaml.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

func (d *Data) HandleDisplay(outputType, outputFile string) {
	switch outputType {
	case "json":
		err := d.SaveToFile(outputFile, "json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Report generated and saved at %s\n", outputFile)
	case "terminal":
		fmt.Println(d.ToYAML())
	case "yml":

		err := d.SaveToFile(outputFile, "yaml")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Report generated and saved at %s\n", outputFile)
	}
}

func (d *Data) SaveToFile(filename, format string) error {

	if filename == "" {
		return fmt.Errorf("filename is required")
	}
	// if file exists  return error
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("file already exists: %s", filename)
	}

	var data string
	var err error

	switch format {
	case "json":
		data, err = d.ToJSON()
	case "yaml":
		data, err = d.ToYAML()
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(data), 0644)
}
