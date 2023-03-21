package util

import (
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type fixedWindow struct {
	lock       sync.Mutex
	lastTime   time.Time
	count      int64
	windowSize int64
	delay      int64
}

func (fixedWindow *fixedWindow) tryAcquire() bool {
	now := time.Now()
	duration := now.Sub(fixedWindow.lastTime).Milliseconds()
	if duration > fixedWindow.delay {
		fixedWindow.count = 0
		fixedWindow.lock.Lock()
		fixedWindow.lastTime = now
		fixedWindow.lock.Unlock()
		return true
	} else {
		fixedWindow.lock.Lock()
		if fixedWindow.count+1 > fixedWindow.windowSize {
			fixedWindow.count = 0
			fixedWindow.lastTime = now
			fixedWindow.lock.Unlock()
			return false
		} else {
			fixedWindow.count++
			fixedWindow.lock.Unlock()
			return true
		}
	}
}

func GetFixedWindow() gin.HandlerFunc {
	window := &fixedWindow{
		lastTime:   time.Now(),
		count:      0,
		windowSize: 10,
		delay:      1000,
	}
	return func(context *gin.Context) {
		if !window.tryAcquire() {
			context.JSON(429, gin.H{
				"message": "too many requests",
			})
			context.Abort()
		}
	}
}
