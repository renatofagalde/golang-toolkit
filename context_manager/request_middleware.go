package contextManager

import (
	headervalidators "bootstrap/src/config/context_manager/header_validators"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestIDMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		errors := make(chan string, 2)
		go headervalidators.ValidateID(c, errors)
		go headervalidators.ValidateJourney(c, errors)

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

		Set(c.Request.Context())

		c.Next()

	}
}
