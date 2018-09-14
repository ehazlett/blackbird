package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/ehazlett/blackbird"
	"github.com/ehazlett/blackbird/server"
	"github.com/ehazlett/blackbird/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.BuildVersion()
	app.Author = "@ehazlett"
	app.Email = ""
	app.Usage = version.Description
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug logging",
		},
		cli.StringFlag{
			Name:  "datastore, d",
			Usage: "datastore backend",
			Value: "memory://",
		},
		cli.StringFlag{
			Name:  "grpc-addr, g",
			Usage: "grpc listen address",
			Value: "127.0.0.1:9000",
		},
		cli.IntFlag{
			Name:  "http-port",
			Usage: "http port",
			Value: 80,
		},
		cli.IntFlag{
			Name:  "https-port",
			Usage: "https port",
			Value: 443,
		},
	}
	app.Action = start
	app.Before = func(ctx *cli.Context) error {
		if ctx.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(ctx *cli.Context) error {
	cfg := &blackbird.Config{
		GRPCAddr:     ctx.String("grpc-addr"),
		DatastoreUri: ctx.String("datastore"),
		HTTPPort:     ctx.Int("http-port"),
		HTTPSPort:    ctx.Int("https-port"),
		Debug:        ctx.Bool("debug"),
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	if err := srv.Run(); err != nil {
		return err
	}

	// wait
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		done <- true
	}()

	<-done
	return nil
}