package url

import (
	"errors"
	"github.com/DaDvoy/url-shortener-api.git/internal/lib/random"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type URLSaver interface {
	SaveURL(urlSave, alias string) error
	GetAlias(url string) (string, error)
}

type Request struct {
	URL string `json:"url"`
}

func (u *Urls) PostURL(c *gin.Context) {
	const op = "handlers.url.PostURL"

	u.Log = u.Log.With(
		slog.String("op", op),
		slog.String("request_id", u.ReqID.GetRequestId()),
	)

	var req Request

	err := c.BindJSON(&req)
	if errors.Is(err, io.EOF) {
		u.Log.Error("request body is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body is empty"})
		return
	}
	if err != nil {
		u.Log.Error("failed to bind json", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if req.URL == "" {
		u.Log.Error("empty url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is missing"})
		return
	}

	alias, err := u.URLSaver.GetAlias(req.URL)
	if errors.Is(err, storage.ErrURLNotFound) {
		alias = random.New()
		err = u.URLSaver.SaveURL(req.URL, alias)
		if err != nil {
			u.Log.Error("failed to save a new alias")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		u.Log.Info("url added")
		c.JSON(http.StatusCreated, gin.H{"short-url": alias})
		return
	}
	if err != nil {
		u.Log.Error("failed to get alias")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	u.Log.Info("url already exists", slog.String("url", req.URL))
	c.JSON(http.StatusFound, gin.H{"short-url": alias})
}
