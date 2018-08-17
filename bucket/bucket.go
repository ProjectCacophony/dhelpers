package bucket

import (
	"math"

	"context"

	"golang.org/x/time/rate"
)

// Bucket is a ratelimit bucket
type Bucket struct {
	limiter *rate.Limiter
}

// NewBucket creates a bucket with the following conditions:
// initially full, refilled at limit tokens per second (limit events per seconds)
func NewBucket(perSecond float64) *Bucket {
	return &Bucket{
		limiter: rate.NewLimiter(
			rate.Limit(perSecond),     // limit
			int(math.Ceil(perSecond)), // set burst to (ceil value of) the rate, allow maximum bursts
		),
	}
}

// Allow returns true when the event may happen now
func (b *Bucket) Allow() bool {
	return b.limiter.Allow()
}

// Wait blocks until event may happen, returns true when the event may happen
// returns false if the context is cancelled or the deadline is exceeded
func (b *Bucket) Wait(ctx context.Context) bool {
	return b.limiter.Wait(ctx) == nil
}
