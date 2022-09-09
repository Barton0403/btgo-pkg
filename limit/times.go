package limit

import (
	"sync"
	"time"
)

type timeLimiter interface {
	Allow() bool
}

type Times struct {
	SecondLimiter timeLimiter
	MinuteLimiter timeLimiter
}

// secondLimiter 一秒内限制
type secondLimiter struct {
	mu          sync.Mutex
	max         int
	count       int
	currentTime int64
}

func (p *secondLimiter) Allow() bool {
	now := time.Now().Unix()

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.currentTime != now {
		p.count = 1
		p.currentTime = now
		return true
	}

	p.count++
	if p.count > p.max {
		return false
	}

	return true
}

// minuteLimiter 一分钟限制
type minuteLimiter struct {
	mu          sync.Mutex
	max         int
	count       int
	currentTime int64
}

func (p *minuteLimiter) Allow() bool {
	now := time.Now().Unix()

	p.mu.Lock()
	defer p.mu.Unlock()

	if now-p.currentTime > 60 {
		p.count = 1
		p.currentTime = now
		return true
	}

	p.count++
	if p.count > p.max {
		return false
	}

	return true
}

func NewTimes(second int, minute int) *Times {
	return &Times{
		SecondLimiter: &secondLimiter{
			max: second,
		},
		MinuteLimiter: &minuteLimiter{
			max: minute,
		},
	}
}
