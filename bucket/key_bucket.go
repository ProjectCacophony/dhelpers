package bucket

import (
	"context"
	"math"

	"sync"

	"golang.org/x/time/rate"
)

// Bucket is a ratelimit bucket with keys
type KeyBucket struct {
	rate  rate.Limit
	burst int

	limiters map[string]*rate.Limiter
	sync.RWMutex
}

// NewKeyBucket creates a bucket with the following conditions:
// initially full, refilled at limit tokens per second (limit events per seconds)
func NewKeyBucket(perSecond float64) *KeyBucket {
	return &KeyBucket{
		rate:     rate.Limit(perSecond),     // limit
		burst:    int(math.Ceil(perSecond)), // set burst to (ceil value of) the rate, allow maximum bursts
		limiters: make(map[string]*rate.Limiter),
	}
}

// Allow returns true when the event may happen now
func (b *KeyBucket) Allow(key string) bool {
	b.Lock()
	defer b.Unlock()

	// reuse existing bucket if possible
	_, ok := b.limiters[key]
	if ok && b.limiters[key] != nil {
		return b.limiters[key].Allow()
	}

	// create new bucket
	b.limiters[key] = rate.NewLimiter(b.rate, b.burst)
	return b.limiters[key].Allow()
}

// Wait blocks until event may happen, returns true when the event may happen
// returns false if the context is cancelled or the deadline is exceeded
func (b *KeyBucket) Wait(ctx context.Context, key string) bool {
	b.Lock()
	defer b.Unlock()

	// reuse existing bucket if possible
	_, ok := b.limiters[key]
	if ok && b.limiters[key] != nil {
		return b.limiters[key].Wait(ctx) == nil
	}

	// create new bucket
	b.limiters[key] = rate.NewLimiter(b.rate, b.burst)
	return b.limiters[key].Wait(ctx) == nil
}
