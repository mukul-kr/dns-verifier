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

var (
	defaultFlag FlagConfig
)

func InitConfig() {

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("INPUT_TYPE", "file")
	viper.SetDefault("INPUT_FILE", "")
	viper.SetDefault("OUTPUT_TYPE", "file")
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
}

func GetConfig() FlagConfig {
	return defaultFlag
}
