package context_manager

import (
	"context"
	"github.com/gin-gonic/gin"
)

func Get() context.Context {
	return context.Background()
}

func GetGinContext() *gin.Context {
	ctx := Get()
	if ctx == nil {
		return nil
	}

	ginCtx, ok := ctx.Value(CTX_KEY).(*gin.Context)
	if !ok {
		return nil
	}
	return ginCtx
}
