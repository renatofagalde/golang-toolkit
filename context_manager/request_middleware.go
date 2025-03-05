package context_manager

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/renatofagalde/golang-toolkit/context_manager/header_validators"
	"net/http"
)

const (
	X_REQUEST_ID      = "X-Request-ID"
	X_REQUEST_JOURNEY = "X-Request-Journey"
	CTX_KEY           = "CTX_KEY"
)

func RequestMiddlewareContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		errors := make(chan string, 2)
		go header_validators.ValidateID(c, errors)
		go header_validators.ValidateJourney(c, errors)

		var errorList []string
		for i := 0; i < 2; i++ {
			if err := <-errors; err != "" {
				errorList = append(errorList, err)
			}
		}

		if len(errorList) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errorList,
			})
			c.Abort()
			return
		}

		// Armazena o contexto da requisição no próprio gin.Context
		ctx := context.WithValue(c.Request.Context(), CTX_KEY, c)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
