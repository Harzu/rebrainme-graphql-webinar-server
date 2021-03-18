package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"

	"webinar/graphql/server/internal/constants"
)

func EnrichContextWithSession(c *gin.Context) {
	sessionToken := c.Request.Header.Get(constants.SessionHeader)
	if sessionToken != "" {
		ctx := context.WithValue(c.Request.Context(), constants.SessionHeader, sessionToken)
		c.Request = c.Request.WithContext(ctx)
	}

	c.Next()
}
