package cmd

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/mukul-kr/dns-verifier/internal/config"
)

func stringValidator(input string) error {

	if len(input) == 0 {
		return errors.New("please enter a valid string path")
	}
	return nil
}

func inputTypeSelector() (string, error) {

	type Option struct {
		Name string
		Ext  string
	}

	options := []Option{
		{Name: "Text File", Ext: "txt"},
		{Name: "JSON File", Ext: "json"},
		{Name: "CSV File", Ext: "csv"},
		{Name: "Terminal", Ext: "terminal"},
	}

	items := make([]string, len(options))
	for i, option := range options {
		items[i] = option.Name
	}

	prompt := promptui.Select{
		Label: "Choose Input Type",
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return options[index].Ext, nil
}

func inputFilePath() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Enter Input File",
		Validate: stringValidator,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatal("Error", "Error", err)
		return "", err
	}

	return result, nil
}

func terminalInput() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter URL",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func outputTypeSelector() (string, error) {
	prompt := promptui.Select{
		Label: "Choose Output Type",
		Items: []string{"YML File", "JSON File", "Terminal"},
	}

	_, result, err := prompt.Run()

	return result, err
}

func outputFilePath() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Enter Output File",
		Validate: stringValidator,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func integerValidator(input string) error {
	if len(input) == 0 {
		return nil
	}

	_, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return errors.New("please enter a valid number")
	}
	return nil
}

func timeoutSelector() (int, error) {
	prompt := promptui.Prompt{
		Label:    "Enter Timeout",
		Validate: integerValidator,
	}
	result, err := prompt.Run()
	if len(result) == 0 {
		cfg := config.GetFlagConfig()
		return cfg.Timeout, nil
	}
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}

func triesSelector() (int, error) {
	prompt := promptui.Prompt{
		Label:    "Enter Tries",
		Validate: integerValidator,
	}

	result, err := prompt.Run()

	if len(result) == 0 {
		cfg := config.GetFlagConfig()
		return cfg.Tries, nil
	}

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}
