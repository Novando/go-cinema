package main

import (
	"github.com/novando/go-cinema/pkg/db/pg"
	"github.com/novando/go-cinema/pkg/env"
	"github.com/novando/go-cinema/pkg/logger"
)

func main() {
	logger.InitZerolog(logger.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		CallerSkip:            3,
		Directory:             "./log",
		Filename:              "logfile",
	})

	// Environment configuration
	if err := env.InitViper("./config/config.local.json"); err != nil {
		panic(err.Error())
	}

	pg.InitPGXv5()
}
