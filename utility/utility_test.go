package utility

import (
	"testing"
)

func TestWithIntegerString(t *testing.T) {
	path := "/users/2"
	base := "/users/"
	integer := 2
	var result int
	result, err := ExtractID(path, base)
	if integer != result || err != nil {
		t.Fatalf(`GetIDFromSting("2") = %d, %v, want match for %d, nil`, result, err, integer)
	}
}

func TestWithInvalidString(t *testing.T) {
	path := "/users/A"
	base := "/users/"
	_, err := ExtractID(path, base)
	if err == nil {
		t.Fatalf(`Should have thrown an error`)
	}
}
