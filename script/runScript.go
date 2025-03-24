package script

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type scriptOptions struct {
	OpList         []interface{}
	WorkingDir     string
	GetOutput      bool
	ParamTransform func(*string)
}

// Most of the content of this function
// comes from here: https://github.com/hashicorp/terraform-provider-external/blob/main/internal/provider/data_source.go
func runScript(o *scriptOptions) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	opList := o.OpList
	workingDir := o.WorkingDir

	if err := validateProgramAttr(opList); err != nil {
		l("Attribute validation failed")
		return "", diag.FromErr(err)
	}

	program := make([]string, len(opList))

	for i, vI := range opList {
		program[i] = vI.(string)
		o.ParamTransform(&program[i])
	}

	l("Command attributes")
	lf(program)

	cmd := exec.Command(program[0], program[1:]...)
	cmd.Dir = workingDir

	if o.GetOutput {
		l("Reading output...")
		resultBytes, err := cmd.Output()
		resultJSON := string(resultBytes)
		log.Printf("[TRACE] JSON output: %+v\r\n", resultJSON)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.Stderr != nil && len(exitErr.Stderr) > 0 {
					return "", diag.Errorf("failed to execute %q: %s", program[0], string(exitErr.Stderr))
				}
				return "", diag.Errorf("command %q failed with no error message", program[0])
			} else {
				return "", diag.Errorf("failed to execute %q: %s", program[0], err)
			}
		}
		return resultJSON, diags
	}

	l("Running script...")
	lf(cmd)
	if resultBytes, err := cmd.Output(); err != nil {
		l("Script returned an error")
		l("OUTPUT")
		l(string(resultBytes))
		var errbuf bytes.Buffer
		cmd.Stderr = &errbuf
		stderr := errbuf.String()
		lf(stderr)
		return "", diag.FromErr(err)
	}
	return "", diags
}
