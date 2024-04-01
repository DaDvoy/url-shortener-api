package storage

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
)

type Storage interface {
	SaveURL(ctx context.Context, urlSave, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
	GetAlias(ctx context.Context, url string) (string, error)
}
