package input

import "testing"

func TestStripArgWithMatch(t *testing.T) {
	arg := "${{ inputs.config-file }}"
	expected := "config-file"
	result, matched := StripArg(arg)

	if !matched {
		t.Errorf("Expected %s to match regex.", arg)
	}

	if result != expected {
		t.Errorf("Expected %s but instead got %s", expected, result)
	}
}

func TestStripArgWithNoMatch(t *testing.T) {
	arg := "${{ env.config-file }}"
	expected := ""
	result, matched := StripArg(arg)

	if matched {
		t.Errorf("Expected %s to match regex.", arg)
	}

	if result != expected {
		t.Errorf("Expected %s but instead got %s", expected, result)
	}
}
