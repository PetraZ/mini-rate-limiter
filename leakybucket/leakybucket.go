package leakybucket

import (
	"sync"
	"time"
)

type Limiter struct {
	*sync.Mutex
	// how many request per second, can be float number here 0.1 mean 1 request per sec
	Rate              float64
	LastUnixTimestamp int64
}

func NewLimiter(rate float64) *Limiter {
	return &Limiter{
		Rate:  rate,
		Mutex: &sync.Mutex{},
	}
}

func (l *Limiter) Take() int64 {
	l.Lock()
	defer l.Unlock()
	now := time.Now().UnixNano()
	nanoSecsSinceLastRequest := float64(now - l.LastUnixTimestamp)
	nanoSecsPerRequest := float64(time.Second.Nanoseconds()) / l.Rate
	if nanoSecsSinceLastRequest >= nanoSecsPerRequest {
		l.LastUnixTimestamp = now
		return now
	}
	time.Sleep(time.Nanosecond * time.Duration(nanoSecsPerRequest-nanoSecsSinceLastRequest))
	now = now + int64(nanoSecsPerRequest-nanoSecsSinceLastRequest)
	l.LastUnixTimestamp = now
	return now
}
