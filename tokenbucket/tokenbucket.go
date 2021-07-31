package tokenbucket

import (
	"sync"
	"time"
)

type Limiter struct {
	*sync.Mutex
	// how many request per second, can be float number here 0.1 mean 1 request per sec
	Rate                    float64
	MaxCapacity             float64
	CurrrentCapacity        float64
	LastRefillUnixTimestamp int64
}

func NewLimiter(rate float64, maxCap float64) *Limiter {
	return &Limiter{
		Rate:        rate,
		MaxCapacity: maxCap,
		Mutex:       &sync.Mutex{},
	}
}

// this needs to be concurrently safe
func (l *Limiter) AllowRequest(consumeCapacity float64) bool {
	l.Lock()
	defer l.Unlock()
	l.refill()

	if l.CurrrentCapacity >= consumeCapacity {
		l.CurrrentCapacity -= consumeCapacity
		return true
	}

	return false
}

func (l *Limiter) refill() {
	now := time.Now().UnixNano()
	refillValue := float64(now-l.LastRefillUnixTimestamp)*l.Rate/1e9 + l.CurrrentCapacity
	l.CurrrentCapacity = min(refillValue, l.MaxCapacity)
	l.LastRefillUnixTimestamp = now
}

func min(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}
