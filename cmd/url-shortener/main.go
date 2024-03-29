package main

import (
	"github.com/DaDvoy/url-shortener-api.git/internal/config"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/handlers/url"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/middleware"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad() // todo: set CONFIG_PATH for launch
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting url-shortener-api", slog.String("env", cfg.Env))

	storage, err := postgres.New()
	if err != nil {
		log.Error("failed to init storage", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		os.Exit(1)
	}

	router := gin.New()
	reqID := middleware.New()
	router.Use(reqID.RequestIdMiddleware)
	router.Use(gin.Logger(), gin.Recovery())

	url := &url.Urls{Log: log, ReqID: reqID, URLSaver: storage, URLReceiver: storage}
	router.GET("/:alias", url.GetURL)
	router.POST("/url", url.PostURL)

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
