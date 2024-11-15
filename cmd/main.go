package main

import (
	"github.com/novando/go-cinema/internal/reservation"
	"github.com/novando/go-cinema/pkg/db/pg"
	"github.com/novando/go-cinema/pkg/env"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/spf13/viper"
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

	db, err := pg.InitPGXv5(pg.Config{
		Host:    viper.GetString("db.pg.host"),
		Port:    viper.GetUint("db.pg.port"),
		User:    viper.GetString("db.pg.user"),
		Pass:    viper.GetString("db.pg.pass"),
		Name:    viper.GetString("db.pg.name"),
		Schema:  viper.GetString("db.pg.schema"),
		MaxPool: viper.GetUint("db.pg.pool"),
		SSL:     viper.GetBool("db.pg.ssl"),
	})
	if err != nil {
		panic(err.Error())
	}

	reservation.Init(db)
}
