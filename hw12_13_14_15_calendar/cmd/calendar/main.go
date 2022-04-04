package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/app"
	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/config"
	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	myConf := config.NewConfig(configFile)
	log := logger.New(myConf.Logger)
	// realStorage := storage.NewStorage(myConf.Storage)
	// calendar := app.New(log, realStorage)
	calendar := app.New(log, memorystorage.New())
	server := internalhttp.NewServer(log, calendar, myConf.Server)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	log.Info("calendar is running...")
	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()
}
