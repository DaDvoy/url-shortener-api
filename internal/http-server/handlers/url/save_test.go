package url_test

import (
	"bytes"
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

func hugeMapPost() (map[string]string, map[string]string) {
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

func TestPostHandler(t *testing.T) {
	r := SetUpRouter()
	skipLog := skipslog.NewSkipLogger()
	strg := in_memory.New()
	fillStorage(strg)
	ur := url.Urls{Log: skipLog, ReqID: middleware.New(), URLSaver: strg, URLReceiver: strg}
	r.POST("/url", ur.PostURL)
	t.Run("already exist", func(t *testing.T) {
		input := fmt.Sprintf(`{"url": "%s"}`, "http://google.com")

		req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusFound)
	})
	t.Run("new added", func(t *testing.T) {
		input := fmt.Sprintf(`{"url": "%s"}`, "http://yandex.ru")

		req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusCreated)
	})
	t.Run("can not twice add one unique url", func(t *testing.T) {
		input := fmt.Sprintf(`{"url": "%s"}`, "http://yandex.ru")

		req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusFound)
	})
	t.Run("request body is empty", func(t *testing.T) {
		input := fmt.Sprintf(`{}`)

		req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusBadRequest)
	})
	t.Run("empty url", func(t *testing.T) {
		input := fmt.Sprintf(`{"url": "%s"}`, "")

		req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusBadRequest)
	})
}

func TestPostHugeMap(t *testing.T) {
	r := SetUpRouter()
	skipLog := skipslog.NewSkipLogger()
	mp1, mp2 := hugeMapPost()
	strg := &in_memory.Storage{UrlToAlias: mp1, AliasToUrl: mp2}
	ur := url.Urls{Log: skipLog, ReqID: middleware.New(), URLSaver: strg, URLReceiver: strg}
	r.POST("/url", ur.PostURL)
	t.Run("huge requests", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			randStr := random.New()

			input := fmt.Sprintf(`{"url": "%s"}`, randStr)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusCreated)
		}
	})
}
