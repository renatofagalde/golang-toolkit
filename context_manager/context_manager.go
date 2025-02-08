package context_manager

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	globalContext context.Context
	globalGinCtx  *gin.Context
	mutex         sync.RWMutex
)

func Set(ctx context.Context, ginCtx *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	globalContext = ctx
	globalGinCtx = ginCtx
	fmt.Sprint("set mutext")
}

func Get() context.Context {
	mutex.RLock()
	defer mutex.RUnlock()
	fmt.Sprint("unlock mutext")
	return globalContext
}

func GetGinContext() *gin.Context {
	mutex.RLock()
	defer mutex.RUnlock()
	return globalGinCtx
}
