package context_manager

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
