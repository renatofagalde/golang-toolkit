package context_manager

const (
	X_REQUEST_ID      = "X-Request-ID"
	X_REQUEST_JOURNEY = "X-Request-Journey"
)

// 20250305
func Give() (journey, requestID string) {
	ginCtx := GetGinContext()
	if ginCtx == nil {
		return "", ""
	}

	if j, exists := ginCtx.Get(X_REQUEST_JOURNEY); exists {
		journey = j.(string)
	}
	if id, exists := ginCtx.Get(X_REQUEST_ID); exists {
		requestID = id.(string)
	}

	return journey, requestID
}
