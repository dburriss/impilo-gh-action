package main

import (
	"testing"
)

// range of vulnerability
// "*"
// "1.0.0-rc1 - 2.1.8"
// "<0.21.2"
// "<=0.21.1"
// ">=0.1.6"
// ">=4.0.0 <4.1.1"
// "<=2.2.3 || 2.2.5 || 3.0.0 - 3.1.4 || >=3.2.0-alpha.0"
// "3.2.1 || 5.0.0 - 5.0.8"

func TestRangeBlank(t *testing.T) {
	input := ""
	expectedFoundIn := "0.0.0"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := ""
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}

func TestRangeAsterisk(t *testing.T) {
	input := "*"
	expectedFoundIn := "0.0.0"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := ""
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}
