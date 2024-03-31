package main

import (
	"flag"
	"github.com/DaDvoy/url-shortener-api.git/internal/app"
	"github.com/DaDvoy/url-shortener-api.git/internal/config"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	in_memory "github.com/DaDvoy/url-shortener-api.git/internal/storage/in-memory"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage/postgres"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	inMemoryFlag := flag.Bool("in-memory", false, "")
	postgresFlag := flag.Bool("postgres", false, "")
	flag.Parse()

	cfg := config.MustLoad()
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting url-shortener-api", slog.String("env", cfg.Env))

	var iStorage storage.Storage

	switch {
	case *inMemoryFlag:
		iStorage = in_memory.New()
		log.Info("storage is in-memory")
	case *postgresFlag:
		var err error
		iStorage, err = postgres.New()
		if err != nil {
			log.Error("failed to init storage", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			os.Exit(1)
		}
		log.Info("storage is postgres")
	default:
		log.Error("Not enough arguments for launch")
		os.Exit(1)
	}

	application := app.New(log, iStorage, 8080)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("application stopping", slog.String("signal", signal.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")

}
