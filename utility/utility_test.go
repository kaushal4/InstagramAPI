package utility

import (
	"testing"
)

func TestWithIntegerString(t *testing.T) {
	integerString := "2"
	integer := 2
	var result int
	result, err := GetIDFromString(integerString)
	if integer != result || err != nil {
		t.Fatalf(`GetIDFromSting("2") = %d, %v, want match for %d, nil`, result, err, integer)
	}
}

func TestWithInvalidString(t *testing.T) {
	integerString := "A"
	_, err := GetIDFromString(integerString)
	if err == nil {
		t.Fatalf(`Should have thrown an error`)
	}
}
