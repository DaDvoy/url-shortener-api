package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
)

type Storage interface {
	SaveURL(urlSave, alias string) error
	GetURL(alias string) (string, error)
	GetAlias(url string) (string, error)
}
