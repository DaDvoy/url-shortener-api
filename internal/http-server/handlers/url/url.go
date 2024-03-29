package url

import (
	"github.com/DaDvoy/url-shortener-api.git/internal/http-server/middleware"
	"golang.org/x/exp/slog"
)

type Urls struct {
	Log         *slog.Logger
	ReqID       *middleware.RequestId
	URLSaver    URLSaver
	URLReceiver UrlReceiver
}
