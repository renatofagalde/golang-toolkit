package context_manager

import (
	"context"
	"fmt"
)

const (
	X_REQUEST_ID      = "X-Request-ID"
	X_REQUEST_JOURNEY = "X-Request-Journey"
)

func Give() (journey, requestID string) {
	ctx := Get()
	if ctx == nil {
		fmt.Println("Context is nil, defaulting to background context")
		ctx = context.Background()
	}

	// Recupera o journey
	if j, ok := ctx.Value(X_REQUEST_JOURNEY).(string); ok {
		journey = j
	} else {
		journey = ""
	}

	// Recupera o requestID
	if id, ok := ctx.Value(X_REQUEST_ID).(string); ok {
		requestID = id
	} else {
		requestID = ""
	}

	return journey, requestID
}
