package dhelpers

import "testing"

func TestHexToDecimal(t *testing.T) {
	var v int
	v = HexToDecimal("ff63eb")
	if v != 16737259 {
		t.Error("Expected 16737259, got ", v)
	}
	v = HexToDecimal("FF63EB")
	if v != 16737259 {
		t.Error("Expected 16737259, got ", v)
	}
	v = HexToDecimal("#ffffff")
	if v != 16777215 {
		t.Error("Expected 16777215, got ", v)
	}
	v = HexToDecimal("#FFFFFF")
	if v != 16777215 {
		t.Error("Expected 16777215, got ", v)
	}
}

func TestDecimalToHex(t *testing.T) {
	var v string
	v = DecimalToHex(16737259)
	if v != "FF63EB" {
		t.Error("Expected FF63EB, got ", v)
	}
}
