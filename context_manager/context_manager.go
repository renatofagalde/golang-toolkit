package contextManager

import (
	"context"
	"fmt"
	"sync"
)

var (
	globalContext context.Context
	mutex         sync.RWMutex
)

func Set(ctx context.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	globalContext = ctx
	fmt.Sprint("set mutext")
}

func Get() context.Context {
	mutex.RLock()
	defer mutex.RUnlock()
	fmt.Sprint("unlock mutext")
	return globalContext
}
