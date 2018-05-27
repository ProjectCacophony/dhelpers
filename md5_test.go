package dhelpers

import "testing"

func TestGetMD5Hash(t *testing.T) {
	v := GetMD5Hash("foobar")
	if v != "3858f62230ac3c915f300c664312c63f" {
		t.Error("Expected 3858f62230ac3c915f300c664312c63f, got ", v)
	}
	v = GetMD5Hash("")
	if v != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Error("Expected d41d8cd98f00b204e9800998ecf8427e, got ", v)
	}
}
