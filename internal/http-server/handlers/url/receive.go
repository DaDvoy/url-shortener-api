package url

import (
	"context"
	"errors"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type UrlReceiver interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

type Response struct {
	Alias string `json:"alias"`
}

func (u *Urls) GetURL(c *gin.Context) {
	const op = "handlers.url.GetURL"

	u.Log = u.Log.With(
		slog.String("op", op),
		slog.String("request_id", u.ReqID.GetRequestId()),
	)

	var response Response

	response.Alias = c.Param("alias")
	if response.Alias == "" {
		u.Log.Info("receive empty short URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "need a non-empty URL"})
		return
	}
	ctx := context.Background()
	url, err := u.URLReceiver.GetURL(ctx, response.Alias)
	if errors.Is(err, storage.ErrURLNotFound) {
		u.Log.Info("non-existent URL", "alias", response.Alias)
		c.JSON(http.StatusNotFound, gin.H{"error": "There is not URL with such a short URL"})
		return
	}
	if err != nil {
		u.Log.Error("failed to get url",
			slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	u.Log.Info("returned URL")
	c.JSON(http.StatusOK, gin.H{"url": url})
}
