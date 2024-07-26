package config

import (
	"github.com/spf13/viper"
)

type FlagConfig struct {
	InputType  string
	InputFile  string
	OutputType string
	OutputFile string
	Timeout    int
	Tries      int
}

type CobraConfig struct {
	Use   string
	Short string
	Long  string
}

var (
	defaultFlag FlagConfig
	cobraConfig CobraConfig
)

func GetFlagConfig() FlagConfig {
	return defaultFlag
}

func GetCobraConfig() CobraConfig {
	return cobraConfig
}

func InitConfig() {

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("INPUT_TYPE", "")
	viper.SetDefault("INPUT_FILE", "")
	viper.SetDefault("OUTPUT_TYPE", "")
	viper.SetDefault("OUTPUT_FILE", "")
	viper.SetDefault("TIMEOUT", 5)
	viper.SetDefault("TRIES", 1)

	defaultFlag = FlagConfig{
		InputType:  viper.GetString("INPUT_TYPE"),
		InputFile:  viper.GetString("INPUT_FILE"),
		OutputType: viper.GetString("OUTPUT_TYPE"),
		OutputFile: viper.GetString("OUTPUT_FILE"),
		Timeout:    viper.GetInt("TIMEOUT"),
		Tries:      viper.GetInt("TRIES"),
	}
	cobraConfig = CobraConfig{
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
	}
}
