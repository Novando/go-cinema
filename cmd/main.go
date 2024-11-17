package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/novando/go-cinema/internal/reservation"
	"github.com/novando/go-cinema/pkg/db/pg"
	"github.com/novando/go-cinema/pkg/env"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/spf13/viper"
	"os"
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

	// init App
	app := fiber.New(fiber.Config{
		AppName:   os.Getenv("APP_NAME"),
		BodyLimit: 8 * 15 * 1024 * 1024,
	})
	v := app.Group("/v1")

	// healtcheck
	v.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint: "/healtz",
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
	}))

	// config cors
	//v.Use(cors.New(cors.Config{
	//	AllowOrigins: strings.Join(cfgParams.CorsList, ","),
	//}))
	v.Use(cors.New())

	reservation.Init(v, db)

	if err = app.Listen(":" + fmt.Sprintf("%v", viper.GetInt("app.port"))); err != nil {
		logger.Call().Fatalf(err.Error())
	}
}
