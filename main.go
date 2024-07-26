package main

import (
	"github.com/mukul-kr/dns-verifier/cmd"
	"github.com/mukul-kr/dns-verifier/internal/config"
	"github.com/mukul-kr/dns-verifier/pkg/logger"
)

func main() {
	config.InitConfig()
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	cmd.Execute(log)
}
