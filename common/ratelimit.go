package common

import (
	"fmt"
	"sync"
	"time"
)

const (
	TimeInverval = 3
)

type RateLimitCtx struct {
	sync.RWMutex
	quotaConfig *SafeMap
	start       time.Time
	quota       []*SafeMap
}

func NewRateLimitCtx(quotaConfig *SafeMap) *RateLimitCtx {
	ctx := &RateLimitCtx{
		quotaConfig: quotaConfig,
		start:       time.Now(),
		quota:       make([]*SafeMap, TimeInverval*10),
	}

	for i := range ctx.quota {
		ctx.quota[i] = NewSafeMap()
	}

	go func() {
		tc := time.Tick(time.Second * TimeInverval)

		for {
			<-tc

			ctx.Lock()
			for i := range ctx.quota {
				if i+TimeInverval < len(ctx.quota) {
					ctx.quota[i] = ctx.quota[i+TimeInverval]
				} else {
					ctx.quota[i] = NewSafeMap()
				}
			}
			ctx.start = time.Now()
			ctx.Unlock()
		}
	}()

	return ctx
}

func (ctx *RateLimitCtx) Acquire(userId int) bool {
	now := time.Now()

	ctx.RLock()
	defer ctx.RUnlock()
	index := now.Sub(ctx.start) / time.Second

	speedI := ctx.quota[index].ReadMap(userId)
	var speed int

	if speedI == nil {
		speed = 0
	} else {
		speed = speedI.(int)
	}

	maxSpeed := ctx.quotaConfig.ReadMap(userId).(int)
	fmt.Println("writeSpeed = ", maxSpeed)
	if speed > maxSpeed {
		fmt.Printf("超过最大读次数了,readSpeed = %d,count = %d ", speed, maxSpeed)
		fmt.Println()
		return false
	}

	ctx.quota[index].WriteMap(userId, speed+1)

	return true
}
