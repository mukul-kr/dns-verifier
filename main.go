package main

import (
	"github.com/mukul-kr/dns-verifier/cmd"
	"github.com/mukul-kr/dns-verifier/pkg/logger"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	cmd.Execute(log)
}
