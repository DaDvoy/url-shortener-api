package in_memory

import "github.com/DaDvoy/url-shortener-api.git/internal/storage"

type Storage struct {
	UrlToAlias map[string]string
	AliasToUrl map[string]string
}

func New() *Storage {
	const op = "storage.in-memory.New"

	mp1 := make(map[string]string)
	mp2 := make(map[string]string)
	s := &Storage{UrlToAlias: mp1, AliasToUrl: mp2}
	return s
}

func (s *Storage) SaveURL(urlSave, alias string) error {
	s.UrlToAlias[urlSave] = alias
	s.AliasToUrl[alias] = urlSave
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	if val, ok := s.AliasToUrl[alias]; ok {
		return val, nil
	}
	return "", storage.ErrURLNotFound
}

func (s *Storage) GetAlias(url string) (string, error) {
	if val, ok := s.UrlToAlias[url]; ok {
		return val, nil
	}
	return "", storage.ErrURLNotFound
}
