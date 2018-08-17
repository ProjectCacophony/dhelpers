package bucket

import (
	"context"
	"testing"
	"time"
)

func TestKeyBucket_Allow(t *testing.T) {
	t.Parallel()

	bucket := NewKeyBucket(2)

	if !bucket.Allow("foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if !bucket.Allow("foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if bucket.Allow("foo") {
		t.Error("expected bucket.Allow() to be false, was true")
	}

	if !bucket.Allow("bar") {
		t.Error("expected bucket.Allow() to be true, was true")
	}

	if !bucket.Allow("bar") {
		t.Error("expected bucket.Allow() to be true, was true")
	}

	if bucket.Allow("bar") {
		t.Error("expected bucket.Allow() to be false, was true")
	}

	time.Sleep(1 * time.Second)

	if !bucket.Allow("foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}
}

func TestKeyBucket_Wait(t *testing.T) {
	t.Parallel()

	bucket := NewKeyBucket(2)

	if !bucket.Wait(context.Background(), "foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if !bucket.Wait(context.Background(), "foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if !bucket.Wait(context.Background(), "bar") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	start := time.Now()

	if !bucket.Wait(context.Background(), "foo") {
		t.Error("expected bucket.Allow() to be true, was false")
	}

	if time.Now().Sub(start).Seconds() < 0.5 {
		t.Error("expected wait time for third call to be above 0.5 seconds")
	}
}
