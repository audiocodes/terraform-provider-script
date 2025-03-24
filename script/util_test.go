package script

import (
	"testing"
)

func TestValidateProgramAttr_valid(t *testing.T) {
	// Test with a command that should be available on most systems
	prog := []interface{}{"echo", "hello"}
	err := validateProgramAttr(prog)
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}
}

func TestValidateProgramAttr_empty(t *testing.T) {
	prog := []interface{}{}
	err := validateProgramAttr(prog)
	if err == nil {
		t.Fatal("Expected error for empty program, got none")
	}
}

func TestValidateProgramAttr_nonString(t *testing.T) {
	prog := []interface{}{123, "hello"}
	err := validateProgramAttr(prog)
	if err == nil {
		t.Fatal("Expected error for non-string program element, got none")
	}
}

func TestValidateProgramAttr_notFound(t *testing.T) {
	prog := []interface{}{"this-command-definitely-does-not-exist", "arg"}
	err := validateProgramAttr(prog)
	if err == nil {
		t.Fatal("Expected error for non-existent program, got none")
	}
}

func TestParseOutput_valid(t *testing.T) {
	output := `{"id":"test-id","resource":"test-resource"}`
	model := parseOutput(output)
	
	if model.ID != "test-id" {
		t.Fatalf("Expected ID 'test-id', got '%s'", model.ID)
	}
	
	if model.Resource != "test-resource" {
		t.Fatalf("Expected Resource 'test-resource', got '%s'", model.Resource)
	}
}

func TestParseOutput_invalid(t *testing.T) {
	output := `not-valid-json`
	model := parseOutput(output)
	
	if model.ID != "" {
		t.Fatalf("Expected empty ID, got '%s'", model.ID)
	}
	
	if model.Resource != "" {
		t.Fatalf("Expected empty Resource, got '%s'", model.Resource)
	}
}

func TestParseOutput_partial(t *testing.T) {
	output := `{"id":"test-id"}`
	model := parseOutput(output)
	
	if model.ID != "test-id" {
		t.Fatalf("Expected ID 'test-id', got '%s'", model.ID)
	}
	
	if model.Resource != "" {
		t.Fatalf("Expected empty Resource, got '%s'", model.Resource)
	}
} 
