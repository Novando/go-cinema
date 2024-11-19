package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/novando/go-cinema/internal/reservation"
	"github.com/novando/go-cinema/pkg/db/pg"
	"github.com/novando/go-cinema/pkg/env"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/spf13/viper"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		AppName:   viper.GetString("app.name"),
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

	// init Socket.io
	c := socket.DefaultServerOptions()
	c.SetServeClient(true)
	c.SetPingInterval(1 * time.Minute)
	c.SetPingTimeout(1 * time.Second)
	c.SetMaxHttpBufferSize(1000000)
	c.SetConnectTimeout(5 * time.Second)
	c.SetCors(&types.Cors{
		Origin: "*",
	})
	soc := socket.NewServer(nil, nil)
	app.Get("/socket.io/", adaptor.HTTPHandler(soc.ServeHandler(c)))
	app.Post("/socket.io/", adaptor.HTTPHandler(soc.ServeHandler(c)))
	soc.On("connection", func(clis ...interface{}) {
		cli := clis[0].(*socket.Socket)
		logger.Call().Infof("%v subs socket.io", cli.Handshake().Address)
	})

	// config cors
	//v.Use(cors.New(cors.Config{
	//	AllowOrigins: strings.Join(cfgParams.CorsList, ","),
	//}))
	v.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	reservation.Init(v, soc, db)

	go app.Listen(":" + fmt.Sprintf("%v", viper.GetInt("app.port")))

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	soc.Close(nil)
	os.Exit(0)
}
