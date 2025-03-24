package script

import (
	"strings"
	"testing"
)

func TestRunScript_validateFailure(t *testing.T) {
	opts := &scriptOptions{
		OpList:     []interface{}{}, // Empty program list (invalid)
		WorkingDir: ".",
		GetOutput:  true,
		ParamTransform: func(s *string) {
			// No transformation
		},
	}

	output, diags := runScript(opts)
	
	if output != "" {
		t.Fatalf("Expected empty output, got: %s", output)
	}
	
	if !diags.HasError() {
		t.Fatal("Expected diagnostics to have errors")
	}
}

func TestRunScript_validCommand(t *testing.T) {
	// Echo a valid JSON that matches our expected format
	opts := &scriptOptions{
		OpList:     []interface{}{"echo", `{"id":"test-id","resource":"test-resource"}`},
		WorkingDir: ".",
		GetOutput:  true,
		ParamTransform: func(s *string) {
			// No transformation
		},
	}

	output, diags := runScript(opts)
	
	if diags.HasError() {
		t.Fatalf("Expected no errors, got: %v", diags)
	}
	
	// Echo adds a newline, need to check if the output contains our expected JSON
	if output == "" || (len(output) > 0 && !strings.Contains(output, `{"id":"test-id","resource":"test-resource"}`)) {
		t.Fatalf("Expected output to contain '{\"id\":\"test-id\",\"resource\":\"test-resource\"}', got '%s'", output)
	}
}

func TestRunScript_paramTransform(t *testing.T) {
	transformed := false
	
	opts := &scriptOptions{
		OpList:     []interface{}{"echo", "##PLACEHOLDER##"},
		WorkingDir: ".",
		GetOutput:  true,
		ParamTransform: func(s *string) {
			if *s == "##PLACEHOLDER##" {
				*s = "transformed"
				transformed = true
			}
		},
	}

	_, diags := runScript(opts)
	
	if diags.HasError() {
		t.Fatalf("Expected no errors, got: %v", diags)
	}
	
	if !transformed {
		t.Fatal("Expected parameter transformation to be applied")
	}
}

func TestRunScript_noOutput(t *testing.T) {
	opts := &scriptOptions{
		OpList:     []interface{}{"echo", "hello"},
		WorkingDir: ".",
		GetOutput:  false,
		ParamTransform: func(s *string) {
			// No transformation
		},
	}

	output, diags := runScript(opts)
	
	if diags.HasError() {
		t.Fatalf("Expected no errors, got: %v", diags)
	}
	
	if output != "" {
		t.Fatalf("Expected empty output when GetOutput is false, got: %s", output)
	}
} 
