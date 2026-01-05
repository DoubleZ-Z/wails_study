package frontend

import (
	"context"
	"sync"
	"wails_study/project/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	mu       sync.RWMutex
	mainCtx  context.Context
	initOnce sync.Once
)

func SetMainContext(ctx context.Context) {
	initOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		mainCtx = ctx
		logger.Info("[Frontend] Main context set")
	})
}

func Emit(event string, payload ...interface{}) {
	mu.RLock()
	ctx := mainCtx
	mu.RUnlock()

	if ctx == nil {
		logger.Warnf("[Frontend] Warning: Emit('%s') ignored — context not set", event)
		return
	}

	select {
	case <-ctx.Done():
		logger.Warnf("[Frontend] Warning: Emit('%s') skipped — context done", event)
		return
	default:
		runtime.EventsEmit(ctx, event, payload...)
	}
}
