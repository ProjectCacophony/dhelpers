package bucket

import (
	"context"
	"testing"
	"time"
)

func TestBucket_Allow(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(2)

	if !bucket.Allow() {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if !bucket.Allow() {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if bucket.Allow() {
		t.Error("expected bucket.Allow() to be false, was true")
	}

	time.Sleep(1 * time.Second)

	if !bucket.Allow() {
		t.Error("expected bucket.Allow() to be true, was false")
	}
}

func TestBucket_Wait(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(2)

	if !bucket.Wait(context.Background()) {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if !bucket.Wait(context.Background()) {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	start := time.Now()

	if !bucket.Wait(context.Background()) {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if time.Since(start).Seconds() < 0.45 {
		t.Error("expected wait time for third call to be above 0.45 seconds")
	}
}
