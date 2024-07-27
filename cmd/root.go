package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mukul-kr/dns-verifier/internal/checker"
	"github.com/mukul-kr/dns-verifier/internal/config"
	"github.com/mukul-kr/dns-verifier/pkg/logger"
	"github.com/mukul-kr/dns-verifier/pkg/reader"
	"github.com/mukul-kr/dns-verifier/pkg/report"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	inputType               string
	inputFile               string
	input                   string
	outputType              string
	outputFile              string
	timeout                 int
	tries                   int
	isMainPersistentFlagSet bool
	log                     *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:   "dns_verifier",
	Short: "A simple DNS configuration for domains are properly set or not",
	Long: `A simple DNS configuration for domains are properly set or not
Text File:
	url1.com,url2.com,...

JSON File:
	[
		{"url": "url1.com"},
		{"url": "url2.com"},
		...
	]

CSV File:
	url
	url1.com
	url2.com
	...`,

	Run: runDNSVerifier,
}

func runDNSVerifier(cmd *cobra.Command, args []string) {
	cfg := config.GetFlagConfig()
	getInputConfig(cfg)
	getOutputConfig(cfg)

	domains, err := reader.HandlerFactory(inputType).Handle(input)()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	data := &report.Data{}

	var wg sync.WaitGroup
	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			data.AddDomain(domain)
			checker.CheckRecords(domain, data)
		}(domain)
	}
	wg.Wait()

	checker.ProcessWg.Wait()

	data.HandleDisplay(outputType, outputFile) // output the data in the output type format

}

// getInputConfig handles input configuration
func getInputConfig(cfg config.FlagConfig) {
	var err error
	if inputType == cfg.InputType && inputFile == cfg.InputFile {
		// log.Info("Using default input values")
		inputType, err = inputTypeSelector()

		log.Info(inputType)
		if err != nil {
			log.Fatal(err)
		}
		if inputType == "terminal" {
			input, err = terminalInput()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			inputFile, err = inputFilePath()

			if err != nil {
				log.Fatal(err)
			}

			typeValidation(inputFile, inputType, "input")

			// read the file
			input = readFile(inputFile)
		}
	} else {
		if inputType == "terminal" {
			input, err = terminalInput()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			typeValidation(inputFile, inputType, "input")

			// read the file
			input = readFile(inputFile)
		}
	}
}

func typeValidation(t1, t2, t string) {
	if !strings.Contains(t1, t2) {
		log.Error(fmt.Sprintf("%s file type and %s type are not same. %s type selected is %s\n", t, t, t, t2))
		os.Exit(1)
	}
}

func readFile(inputFile string) string {
	// read the file, remember to handle io error correctly
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(file)
}

// getOutputConfig handles output configuration
func getOutputConfig(cfg config.FlagConfig) {
	var err error
	if outputType == cfg.OutputType && outputFile == cfg.OutputFile {
		// log.Info("Using default output values")
		outputType, err = outputTypeSelector()
		if err != nil {
			log.Fatal(err)
		}
		if outputType != "terminal" {
			outputFile, err = outputFilePath()
			if err != nil {
				log.Fatal(err)
			}
			typeValidation(outputFile, outputType, "output")
		}
	}
}

func Execute(l *logger.Logger) {
	log = l.Named("cmd")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	initializeFlags()
}

func initializeFlags() {

	cfg := config.GetFlagConfig()

	rootCmd.PersistentFlags().StringVarP(&inputType, "input-type", "i", cfg.InputType, "Input type for the DNS Verifier (csv, json, txt, terminal)")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "f", cfg.InputFile, "Input file path for the DNS Verifier( Not needed for terminal input )")
	rootCmd.PersistentFlags().StringVarP(&outputType, "output-type", "o", cfg.OutputType, "Output type for the DNS Verifier (json, yml, terminal)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "O", cfg.OutputFile, "Output file for the DNS Verifier ( Not needed for terminal output )")

	isMainPersistentFlagSet = len(inputFile)+len(inputType)+len(outputFile)+len(outputType) > 0
}
