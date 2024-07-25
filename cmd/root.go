package cmd

import (
	"os"

	"github.com/mukul-kr/dns-verifier/internal/config"
	"github.com/mukul-kr/dns-verifier/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	inputType  string
	inputFile  string
	input      string
	outputType string
	outputFile string
	timeout    int
	tries      int
	log        *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:   "dns_verifier",
	Short: "A simple DNS configuration for domains are properly set or not",
	Long: `A simple DNS configuration for domains are properly set or not
Text File:
	url1.com
	url2.com
	...

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

	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		if inputType == cfg.InputType && inputFile == cfg.InputFile {
			log.Info("Using default values")
			res, err := inputTypeSelector()
			if err != nil {
				log.Fatal(err)
			}
			inputType = res
			if inputType == "Terminal" {
				res, err := terminalInput()
				if err != nil {
					log.Fatal(err)
				}
				input = res
			} else {
				res, err := inputFilePath()
				if err != nil {
					log.Fatal(err)
				}
				inputFile = res

				// read file
			}
		}

		if outputType == cfg.OutputType && outputFile == cfg.OutputFile {
			log.Info("Using default values")
			res, err := outputTypeSelector()
			if err != nil {
				log.Fatal(err)
			}
			outputType = res
			if outputType != "Terminal" {
				res, err := outputFilePath()
				if err != nil {
					log.Fatal(err)
				}
				outputFile = res
			}
		}
		timeout, err := timeoutSelector()
		if err != nil {
			log.Fatal(err)
		}
		tries, err := triesSelector()
		if err != nil {
			log.Fatal(err)
		}

		log.Infow("Configuration", "Input Type", inputType, "Input File", inputFile, "Output Type", outputType, "Output File", outputFile, "Timeout", timeout, "Tries", tries)

	},
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

	cfg := config.GetConfig()

	rootCmd.PersistentFlags().StringVarP(&inputType, "input-type", "i", cfg.InputType, "Input type for the DNS Verifier (csv, json, txt)")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "f", cfg.InputFile, "Input file for the DNS Verifier")
	rootCmd.PersistentFlags().StringVarP(&outputType, "output-type", "o", cfg.OutputType, "Output type for the DNS Verifier")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "O", cfg.OutputFile, "Output file for the DNS Verifier")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", cfg.Timeout, "Timeout for the network calls made( in seconds )")
	rootCmd.PersistentFlags().IntVarP(&tries, "tries", "T", cfg.Tries, "Tries for the DNS Verifier")

}
