package script

//original source: https://github.com/hashicorp/terraform-provider-external/blob/2b1150d04771816bae85ca0d162f9c8e12c6c52a/internal/provider/util.go#L12

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// validateProgramAttr is a validation function for the "program" attribute we
// accept as input on our resources.
//
// The attribute is assumed to be specified in schema as a list of strings.
func validateProgramAttr(v interface{}) error {
	args := v.([]interface{})
	if len(args) < 1 {
		return fmt.Errorf("'program' list must contain at least one element")
	}

	for i, vI := range args {
		if _, ok := vI.(string); !ok {
			return fmt.Errorf(
				"'program' element %d is %T; a string is required",
				i, vI,
			)
		}
	}

	// first element is assumed to be an executable command, possibly found
	// using the PATH environment variable.
	_, err := exec.LookPath(args[0].(string))
	if err != nil {
		return fmt.Errorf("can't find external program %q", args[0])
	}

	return nil
}
func l(msg string) {
	log.Printf("[TRACE] %s", msg)
}
func lf(v ...interface{}) {
	log.Printf("[TRACE] %v", v)
}
func parseOutput(output string) *scriptModel {
	model := &scriptModel{}
	json.Unmarshal([]byte(output), model)
	return model
}
