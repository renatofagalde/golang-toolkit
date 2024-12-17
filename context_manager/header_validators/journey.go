package header_validators

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

const X_REQUEST_JOURNEY = "X-Request-Journey"

func ValidateJourney(c *Context, errors chan string) {
	journey := c.GetHeader(X_REQUEST_JOURNEY)
	if journey == "" {
		errors <- fmt.Sprintf("Header %s not found", X_REQUEST_JOURNEY)
	} else {
		ctx := context.WithValue(c.Request.Context(), X_REQUEST_JOURNEY, journey)
		c.Request = c.Request.WithContext(ctx)
		errors <- ""
	}

}
