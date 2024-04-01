package services_test

import (
	"context"
	"fmt"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/logger/skipslog"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/random"
	"github.com/DaDvoy/url-shortener-api.git/internal/services"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	in_memory "github.com/DaDvoy/url-shortener-api.git/internal/storage/in-memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func hugeMapSave() (map[string]string, map[string]string) {
	mp1 := make(map[string]string, 100)
	mp2 := make(map[string]string, 100)
	for i := 0; i < 100; i++ {
		rn1 := random.New()
		rn2 := random.New()
		if rn1 != rn2 {
			mp1[rn1] = rn2
			mp2[rn2] = rn1
		}
	}
	return mp1, mp2
}

func hugeMapGet() (map[int]string, map[string]string) {
	mp1 := make(map[int]string, 100)
	mp2 := make(map[string]string, 100)
	for i := 0; i < 100; i++ {
		rn1 := random.New()
		rn2 := random.New()
		if rn1 != rn2 {
			mp1[i] = rn2
			mp2[rn2] = rn1
		}
	}
	return mp1, mp2
}

func getMap() map[string]string {
	return map[string]string{
		"http://google.com": "L7du_MycqX",
		"http://golang.org": "SoacZ1yCDu",
		"http://mysite.ru":  "cJQtkiv0zX",
	}
}

func fillStorage(storage *in_memory.Storage) {
	mp := getMap()
	ctx := context.TODO()
	for v := range mp {
		_ = storage.SaveURL(ctx, v, mp[v])
	}
}

func TestSaveURL(t *testing.T) {
	strg := in_memory.New()
	fillStorage(strg)
	ctx := context.TODO()
	s := &services.Services{Log: skipslog.NewSkipLogger(), URLSaver: strg, UrlReceiver: strg}
	t.Run("already exist", func(t *testing.T) {
		_, err := s.SaveURL(ctx, "http://google.com")
		require.Equal(t, err, storage.ErrURLExists)
	})
	t.Run("added new url", func(t *testing.T) {
		_, err := s.SaveURL(ctx, "http://go.dev")
		require.Equal(t, err, nil)
	})
	t.Run("huge requests", func(t *testing.T) {
		mp1, mp2 := hugeMapSave()
		strg = &in_memory.Storage{UrlToAlias: mp1, AliasToUrl: mp2}
		s = &services.Services{Log: skipslog.NewSkipLogger(), URLSaver: strg, UrlReceiver: strg}
		for i := 0; i < 1000; i++ {
			randStr := random.New()

			_, err := s.SaveURL(ctx, randStr)
			require.NoError(t, err)
			require.Equal(t, err, nil)
		}
	})
}

func TestGetURL(t *testing.T) {
	strg := in_memory.New()
	fillStorage(strg)
	ctx := context.TODO()
	s := &services.Services{Log: skipslog.NewSkipLogger(), URLSaver: strg, UrlReceiver: strg}
	t.Run("find url by short url", func(t *testing.T) {
		mp := getMap()
		for v := range mp {
			shortUrl := fmt.Sprintf("%s", mp[v])
			_, err := s.GetURL(ctx, shortUrl)
			require.NoError(t, err)
		}
	})
	t.Run("find url by short url", func(t *testing.T) {
		_, err := s.GetURL(ctx, "")
		require.Equal(t, err, storage.ErrURLNotFound)
	})
	t.Run("huge map", func(t *testing.T) {
		mp1, mp2 := hugeMapGet()
		strg = &in_memory.Storage{AliasToUrl: mp2}
		s = &services.Services{Log: skipslog.NewSkipLogger(), URLSaver: strg, UrlReceiver: strg}
		for i := 0; i < len(mp2); i++ {
			alias := fmt.Sprintf("%s", mp1[i])
			_, err := s.GetURL(ctx, alias)
			require.NoError(t, err)
		}
	})

}
