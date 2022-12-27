package main

import (
	"testing"
)

// range of vulnerability

// " "
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

// "*"
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

// "1.0.0-rc1 - 2.1.8"
func TestRangeFromTo(t *testing.T) {
	input := "1.0.0-rc1 - 2.1.8"
	expectedFoundIn := "1.0.0-rc1"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := ">2.1.8"
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}

// ">=4.0.0 <4.1.1"
func TestRange(t *testing.T) {
	input := ">=4.0.0 <4.1.1"
	expectedFoundIn := "4.0.0"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := "4.1.1"
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}

// "<=0.21.1"
func TestRangeLessThanOrEqual(t *testing.T) {
	input := "<=0.21.1"
	expectedFoundIn := "0.0.0"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := ">0.21.1"
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}

// "<=0.21.1"
func TestRangeLessThan(t *testing.T) {
	input := "<0.21.1"
	expectedFoundIn := "0.0.0"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := "0.21.1"
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}

// ">=0.1.6"
func TestRangeGreaterThanOrEqual(t *testing.T) {
	input := ">=0.1.6"
	expectedFoundIn := "0.1.6"

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

// ">0.1.6"
func TestRangeGreaterThan(t *testing.T) {
	input := ">0.1.6"
	expectedFoundIn := ">0.1.6"

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

// "<=2.2.3 || 2.2.5 || 3.0.0 - 3.1.4 || >=3.2.0-alpha.0"
func TestRangeMultiple(t *testing.T) {
	input := "<=2.2.3 || 2.2.5 || 3.0.0 - 3.1.4 || >=3.2.0-alpha.0"
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

// "3.2.1 || 5.0.0 - 5.0.8"
func TestMultipleWithRange(t *testing.T) {
	input := "3.2.1 || 5.0.0 - 5.0.8"
	expectedFoundIn := "3.2.1"

	rangeM := parseNpmRanges(input)
	actualFoundIn, exists := rangeM["foundIn"]
	if !exists || actualFoundIn != expectedFoundIn {
		t.Errorf("Expected %s, instead got %s", expectedFoundIn, actualFoundIn)
	}

	expectedFixedIn := ">5.0.8"
	actualFixedIn, exists := rangeM["fixedIn"]
	if !exists || actualFixedIn != expectedFixedIn {
		t.Errorf("Expected %s, instead got %s", expectedFixedIn, actualFixedIn)
	}
}
