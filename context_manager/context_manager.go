package context_manager

import (
	"context"
	"github.com/gin-gonic/gin"
)

func Get() context.Context {
	ginCtx := GetGinContext()
	if ginCtx != nil {
		return ginCtx.Request.Context()
	}
	return context.Background()
}

func GetGinContext() *gin.Context {
	ginCtx, exists := context.WithValue(context.Background(), CTX_KEY, nil).Value(CTX_KEY).(*gin.Context)
	if exists {
		return ginCtx
	}
	return nil
}
