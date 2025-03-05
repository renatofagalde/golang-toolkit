package header_validators

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const X_REQUEST_ID = "X-Request-ID"

func ValidateID(c *gin.Context, errors chan string) {
	requestID := c.GetHeader(X_REQUEST_ID)
	if requestID == "" {
		errors <- fmt.Sprintf("Header %s not found", X_REQUEST_ID)
		return
	}

	c.Set(X_REQUEST_ID, requestID)

	errors <- ""
}
