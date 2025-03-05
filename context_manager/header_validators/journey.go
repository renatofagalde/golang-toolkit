package header_validators

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const X_REQUEST_JOURNEY = "X-Request-Journey"

func ValidateJourney(c *gin.Context, errors chan string) {
	journey := c.GetHeader(X_REQUEST_JOURNEY)
	if journey == "" {
		errors <- fmt.Sprintf("Header %s not found", X_REQUEST_JOURNEY)
		return
	}

	c.Set(X_REQUEST_JOURNEY, journey)

	errors <- ""
}
