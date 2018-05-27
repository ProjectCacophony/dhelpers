package humanize

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	v := Duration(time.Hour*48 + time.Hour*12)
	if v != "2d12h" {
		t.Error("Expected 2d12h, got ", v)
	}
	v = Duration(time.Hour * 24)
	if v != "1d" {
		t.Error("Expected 1d, got ", v)
	}
	v = Duration(time.Hour * 12)
	if v != "12h" {
		t.Error("Expected 12h, got ", v)
	}
	v = Duration(time.Hour*12 + time.Minute*30)
	if v != "12h30m" {
		t.Error("Expected 12h30m, got ", v)
	}
	v = Duration(time.Hour*12 + time.Minute*30 + time.Second*30)
	if v != "12h30m30s" {
		t.Error("Expected 12h30m30s, got ", v)
	}
}
