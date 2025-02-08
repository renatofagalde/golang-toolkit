package context_manager

import (
	"fmt"
)

const (
	X_REQUEST_ID      = "X-Request-ID"
	X_REQUEST_JOURNEY = "X-Request-Journey"
)

func Give() (journey, requestID string) {
	ginCtx := GetGinContext()
	if ginCtx != nil {
		if j, exists := ginCtx.Get(X_REQUEST_JOURNEY); exists {
			journey = j.(string)
		}
		if id, exists := ginCtx.Get(X_REQUEST_ID); exists {
			requestID = id.(string)
		}
	}

	ctx := Get()
	if ctx == nil {
		fmt.Println("Context is nil, defaulting to background context")
		ctx = ginCtx.Request.Context()
	}

	if journey == "" {
		if j, ok := ctx.Value(X_REQUEST_JOURNEY).(string); ok {
			journey = j
		}
	}

	if requestID == "" {
		if id, ok := ctx.Value(X_REQUEST_ID).(string); ok {
			requestID = id
		}
	}

	return journey, requestID
}
