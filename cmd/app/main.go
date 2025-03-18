package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/x0k/skillrock-tasks-service/internal/app"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
)

var (
	config_path string
)

func init() {
	flag.StringVar(&config_path, "config", os.Getenv("CONFIG_PATH"), "Config path")
	flag.Parse()
}

func main() {
	cfg := app.MustLoadConfig(config_path)
	log := app.MustNewLogger(&cfg.Logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info(ctx, "press CTRL-C to exit")
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		s := <-stop
		log.Info(ctx, "signal received", slog.String("signal", s.String()))
		cancel()
	}()

	if err := app.Run(ctx, cfg, log); err != nil {
		log.Error(ctx, "failde shutdown", sl.Err(err))
		os.Exit(1)
	}
	log.Info(ctx, "graceful shutdown")
}
