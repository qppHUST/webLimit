package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type timeSlot struct {
	startTime time.Time
	count     int
}

type SlideWindowLimitRate struct {
	rate           int
	windowDuration time.Duration
	slotDuration   time.Duration

	mu       sync.Mutex
	slotList []*timeSlot
}

func NewSlideWindowLimitRate(rate int, windowDuration time.Duration, slotDuration time.Duration) *SlideWindowLimitRate {
	tool := &SlideWindowLimitRate{
		rate:           rate,
		windowDuration: windowDuration,
		slotDuration:   slotDuration,
		slotList:       make([]*timeSlot, windowDuration/slotDuration),
	}
	now := time.Now()
	delta := time.Second
	for i := range tool.slotList {
		tool.slotList[i] = &timeSlot{
			startTime: now.Add(delta),
			count:     0,
		}
		delta += tool.slotDuration
	}
	return tool
}

func (r *SlideWindowLimitRate) Acquire() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	index := int(now.Sub(r.slotList[0].startTime).Seconds() / r.slotDuration.Seconds())
	if index > len(r.slotList)-1 {
		delta := index - len(r.slotList) + 1
		for i := range r.slotList {
			slot := r.slotList[i]
			for i := 0; i < delta; i++ {
				slot.startTime = slot.startTime.Add(time.Second)
			}
			if i < len(r.slotList)-delta {
				slot.count = r.slotList[i+delta].count
			} else {
				slot.count = 0
			}
		}
		return true
	} else {
		count := 0
		for i := range r.slotList {
			count += r.slotList[i].count
		}
		if count > r.rate {
			return false
		}
		r.slotList[index].count++
		return true
	}
}

func (r *SlideWindowLimitRate) newTimeSlot(startTime time.Time) *timeSlot {
	return &timeSlot{startTime: startTime, count: 1}
}

func GetSlidingWindowHandler() gin.HandlerFunc {
	rate := NewSlideWindowLimitRate(3, time.Second*3, time.Second)
	fmt.Println()
	return func(context *gin.Context) {
		if !rate.Acquire() {
			context.JSON(429, gin.H{
				"message": "too many requests",
			})
		}
		context.Abort()
	}
}
