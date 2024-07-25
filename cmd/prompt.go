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
	prompt := promptui.Select{
		Label: "Choose Input Type",
		Items: []string{"Text File", "JSON File", "CSV File", "Terminal"},
	}

	_, result, err := prompt.Run()

	return result, err
}

func inputFilePath() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Enter Input File",
		Validate: stringValidator,
	}

	result, err := prompt.Run()
	if err != nil {
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
	log.Info(result)
	if len(result) == 0 {
		cfg := config.GetConfig()
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
		cfg := config.GetConfig()
		return cfg.Timeout, nil
	}

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(result)
}
