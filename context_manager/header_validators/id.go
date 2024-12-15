package headervalidators

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

const X_REQUEST_ID = "X-Request-ID"

func ValidateID(c *gin.Context, errors chan string) {
	journey := c.GetHeader(X_REQUEST_ID)
	if journey == "" {
		errors <- fmt.Sprintf("Header %s not found", X_REQUEST_ID)
	} else {
		ctx := context.WithValue(c.Request.Context(), X_REQUEST_ID, journey)
		c.Request = c.Request.WithContext(ctx)
		errors <- ""
	}

}
