package skipslog

import (
	"context"
	"golang.org/x/exp/slog"
)

func NewSkipLogger() *slog.Logger {
	return slog.New(NewSkipHandler())
}

type SkipHandler struct{}

func NewSkipHandler() *SkipHandler {
	return &SkipHandler{}
}

func (s *SkipHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (s *SkipHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return s
}

func (s *SkipHandler) WithGroup(_ string) slog.Handler {
	return s
}

func (s *SkipHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
