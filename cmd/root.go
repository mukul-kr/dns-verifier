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
	Use:   config.GetCobraConfig().Use,
	Short: config.GetCobraConfig().Short,
	Long:  config.GetCobraConfig().Long,

	Run: runDNSVerifier,
}

// runDNSVerifier contains the main logic of the DNS verifier
func runDNSVerifier(cmd *cobra.Command, args []string) {
	cfg := config.GetFlagConfig()
	getInputConfig(cfg)
	getOutputConfig(cfg)
	getTimeoutConfig(cfg)
	getTriesConfig(cfg)

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

	data.HandleDisplay(outputType, outputFile)
	// output the data

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
			// read the file

			typeValidation(inputFile, inputType, "input")

			input = readFile(inputFile)
		}
	} else {
		// read the file
		typeValidation(inputFile, inputType, "input")
		input = readFile(inputFile)
	}
}

func typeValidation(t1, t2, t string) {
	if !strings.Contains(t1, t2) {
		fmt.Println(fmt.Sprintf("%s file type and %s type are not same", t, t))
		os.Exit(1)
	}
}

func readFile(inputFile string) string {
	// read the file
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

// getTimeoutConfig handles timeout configuration
func getTimeoutConfig(cfg config.FlagConfig) {

	if timeout == -1 {
		if isMainPersistentFlagSet {
			timeout = cfg.Timeout
		} else {
			var err error
			timeout, err = timeoutSelector()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// getTriesConfig handles tries configuration
func getTriesConfig(cfg config.FlagConfig) {
	if isMainPersistentFlagSet && tries == -1 {
		tries = cfg.Tries
	} else if tries == -1 {
		var err error
		tries, err = triesSelector()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tries = cfg.Tries
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
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", -1, "Timeout for the network calls made( in seconds )")
	rootCmd.PersistentFlags().IntVarP(&tries, "tries", "T", -1, "Tries for the DNS Verifier")

	isMainPersistentFlagSet = len(inputFile)+len(inputType)+len(outputFile)+len(outputType) > 0
}
