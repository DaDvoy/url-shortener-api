package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/random"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	"golang.org/x/exp/slog"
)

type Services struct {
	Log         *slog.Logger
	URLSaver    URLSaver
	UrlReceiver UrlReceiver
}

type URLSaver interface {
	SaveURL(ctx context.Context, urlSave, alias string) error
	GetAlias(ctx context.Context, url string) (string, error)
}

type UrlReceiver interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, saver URLSaver, receiver UrlReceiver) *Services {
	return &Services{
		Log:         log,
		URLSaver:    saver,
		UrlReceiver: receiver,
	}
}

func (s *Services) SaveURL(ctx context.Context, url string) (string, error) {
	const op = "services.SaveURL"

	log := s.Log.With(
		slog.String("op", op),
		slog.String("url", url),
	)

	alias, err := s.URLSaver.GetAlias(ctx, url)
	if errors.Is(err, storage.ErrURLNotFound) {
		alias = random.New()
		err = s.URLSaver.SaveURL(ctx, url, alias)
		if err != nil {
			log.Error("failed to save a new short url")
			return "", fmt.Errorf("%s: %w", op, errors.New("internal error"))
		}
		log.Info("url added", "short-url", alias)
		return alias, nil
	}
	if err != nil {
		log.Error("failed to get alias")
		return "", fmt.Errorf("%s: %w", op, errors.New("internal error"))
	}

	log.Info("url already exists", "short-url", alias)
	return alias, storage.ErrURLExists
}

func (s *Services) GetURL(ctx context.Context, shorURL string) (string, error) {
	const op = "services.GetURL"

	log := s.Log.With(
		slog.String("op", op),
		slog.String("short-url", shorURL),
	)

	url, err := s.UrlReceiver.GetURL(ctx, shorURL)
	if errors.Is(err, storage.ErrURLNotFound) {
		log.Info("non-existent URL", "alias", shorURL)
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		log.Error("failed to get url",
			slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		return "", fmt.Errorf("%s: %w", op, errors.New("internal error"))
	}
	log.Info("returned URL", "url", url)
	return url, nil
}
