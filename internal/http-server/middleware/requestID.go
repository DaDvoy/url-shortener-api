package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestId struct {
	ReqId string
}

func New() *RequestId {
	return &RequestId{}
}

func (r *RequestId) RequestIdMiddleware(c *gin.Context) {
	r.ReqId = c.GetHeader("X-Request-ID")

	if r.ReqId == "" {
		r.ReqId = uuid.New().String()
		c.Request.Header.Set("X-Request-ID", r.ReqId)
	}
	c.Writer.Header().Set("X-Request-ID", r.ReqId)

	c.Next()
}

func (r *RequestId) GetRequestId() string {
	return r.ReqId
}
