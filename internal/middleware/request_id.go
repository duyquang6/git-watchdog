package middleware

import (
	"net/http"

	"github.com/duyquang6/git-watchdog/pkg/ulid"

	"github.com/duyquang6/git-watchdog/internal/controller"
	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader = "X-Request-ID"
)

// PopulateRequestID add request ID for debugging and tracing request
func PopulateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string
		ctx := c.Request.Context()

		if _, ok := c.Request.Header[http.CanonicalHeaderKey(requestIDHeader)]; ok {
			requestID = c.Request.Header[http.CanonicalHeaderKey(requestIDHeader)][0]
		} else {
			requestID = ulid.GetUniqueID()
		}

		ctx = controller.WithRequestID(ctx, requestID)
		c.Request = c.Request.Clone(ctx)

		c.Next()
	}
}
