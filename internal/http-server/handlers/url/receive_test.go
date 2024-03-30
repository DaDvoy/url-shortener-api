package url_test

import (
	"fmt"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/handlers/url"
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/middleware"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/logger/skipslog"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/random"
	in_memory "github.com/DaDvoy/url-shortener-api.git/internal/storage/in-memory"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestGetHandler(t *testing.T) {
	r := SetUpRouter()
	skipLog := skipslog.NewSkipLogger()
	strg := in_memory.New()
	fillStorage(strg)
	ur := url.Urls{Log: skipLog, ReqID: middleware.New(), URLSaver: strg, URLReceiver: strg}
	r.GET("/:alias", ur.GetURL)
	t.Run("find url by short url", func(t *testing.T) {
		mp := getMap()
		for v := range mp {
			alias := fmt.Sprintf("/%s", mp[v])

			req, err := http.NewRequest("GET", alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)
		}
	})
	t.Run("empty field with short url", func(t *testing.T) {
		req, err := http.NewRequest("GET", "", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusNotFound)
	})
	t.Run("not exist short url", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/dfsfsdf", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusNotFound)
	})
}

func TestGetHugeMap(t *testing.T) {
	r := SetUpRouter()
	skipLog := skipslog.NewSkipLogger()
	mp1, mp2 := hugeMapGet()
	strg := &in_memory.Storage{AliasToUrl: mp2}
	ur := url.Urls{Log: skipLog, ReqID: middleware.New(), URLSaver: strg, URLReceiver: strg}
	r.GET("/:alias", ur.GetURL)
	t.Run("huge map", func(t *testing.T) {
		for i := 0; i < len(mp2); i++ {
			alias := fmt.Sprintf("/%s", mp1[i])

			req, err := http.NewRequest("GET", alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)
		}
	})
}
