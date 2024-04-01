package in_memory

import (
	"context"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
)

type Storage struct {
	UrlToAlias map[string]string
	AliasToUrl map[string]string
}

func New() *Storage {
	mp1 := make(map[string]string)
	mp2 := make(map[string]string)
	s := &Storage{UrlToAlias: mp1, AliasToUrl: mp2}
	return s
}

func (s *Storage) SaveURL(ctx context.Context, urlSave, alias string) error {
	s.UrlToAlias[urlSave] = alias
	s.AliasToUrl[alias] = urlSave
	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	if val, ok := s.AliasToUrl[alias]; ok {
		return val, nil
	}
	return "", storage.ErrURLNotFound
}

func (s *Storage) GetAlias(ctx context.Context, url string) (string, error) {
	if val, ok := s.UrlToAlias[url]; ok {
		return val, nil
	}
	return "", storage.ErrURLNotFound
}
