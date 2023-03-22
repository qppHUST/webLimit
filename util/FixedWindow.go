package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type fixedWindow struct {
	lock       sync.Mutex
	lastTime   time.Time
	count      int
	windowSize int
	delay      time.Duration
}

func (fixedWindow *fixedWindow) TryAcquire() bool {
	fixedWindow.lock.Lock()
	defer fixedWindow.lock.Unlock()
	now := time.Now()
	fmt.Println("当前请求时间为： ", now)
	duration := now.Sub(fixedWindow.lastTime)
	if duration > fixedWindow.delay {
		fmt.Println("超过delay，请求成功")
		fixedWindow.count = 1
		fixedWindow.lastTime = now
		return true
	} else {
		if fixedWindow.count+1 > fixedWindow.windowSize {
			fmt.Println("仍处于duration内部", "count为: ", fixedWindow.count+1)
			fixedWindow.count = fixedWindow.count + 1
			return false
		} else {
			fmt.Println("仍处于duration内部", "count为: ", fixedWindow.count+1)
			fixedWindow.count = fixedWindow.count + 1
			return true
		}
	}
}

func NewFixedWindow(windowSize int, delay time.Duration) *fixedWindow {
	return &fixedWindow{
		lastTime:   time.Now(),
		count:      0,
		windowSize: windowSize,
		delay:      delay,
	}
}

func GetFixedWindowHandler() gin.HandlerFunc {
	window := &fixedWindow{
		lastTime:   time.Now(),
		count:      0,
		windowSize: 10,
		delay:      time.Second * 3,
	}
	return func(context *gin.Context) {
		if !window.TryAcquire() {
			context.JSON(429, gin.H{
				"message": "too many requests",
			})
			context.Abort()
		}
	}
}
