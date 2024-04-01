package main

import (
	"flag"
	"github.com/DaDvoy/url-shortener-api.git/internal/config"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/handlers/url"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/middleware"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	in_memory "github.com/DaDvoy/url-shortener-api.git/internal/storage/in-memory"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
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

	router := gin.New()
	reqID := middleware.New()
	router.Use(reqID.RequestIdMiddleware)
	router.Use(gin.Logger(), gin.Recovery())

	u := &url.Urls{Log: log, ReqID: reqID, URLSaver: iStorage, URLReceiver: iStorage}
	router.GET("/:alias", u.GetURL)
	router.POST("/url", u.PostURL)

	srvListen := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srvListen.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
