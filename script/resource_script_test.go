package script

import (
	"testing"
	"os"
	"path/filepath"
	"io/ioutil"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceScript(t *testing.T) {
	r := resourceScript()
	if r.Schema["read"] == nil {
		t.Fatal("Expected 'read' in schema")
	}
	if r.Schema["create"] == nil {
		t.Fatal("Expected 'create' in schema")
	}
	if r.Schema["update"] == nil {
		t.Fatal("Expected 'update' in schema")
	}
	if r.Schema["delete"] == nil {
		t.Fatal("Expected 'delete' in schema")
	}
	if r.Schema["target_state"] == nil {
		t.Fatal("Expected 'target_state' in schema")
	}
	if r.Schema["working_dir"] == nil {
		t.Fatal("Expected 'working_dir' in schema")
	}
	if r.Schema["resource"] == nil {
		t.Fatal("Expected 'resource' in schema")
	}
}

func TestResourceScript_required(t *testing.T) {
	r := resourceScript()
	
	if !r.Schema["read"].Required {
		t.Fatal("Expected 'read' to be required")
	}
	if !r.Schema["create"].Required {
		t.Fatal("Expected 'create' to be required")
	}
	if !r.Schema["update"].Required {
		t.Fatal("Expected 'update' to be required")
	}
	if !r.Schema["delete"].Required {
		t.Fatal("Expected 'delete' to be required")
	}
	if !r.Schema["target_state"].Required {
		t.Fatal("Expected 'target_state' to be required")
	}
	if !r.Schema["working_dir"].Required {
		t.Fatal("Expected 'working_dir' to be required")
	}
	if r.Schema["resource"].Required {
		t.Fatal("Expected 'resource' to not be required")
	}
	if !r.Schema["resource"].Computed {
		t.Fatal("Expected 'resource' to be computed")
	}
}

func TestResourceScript_types(t *testing.T) {
	r := resourceScript()
	
	if r.Schema["read"].Type != schema.TypeList {
		t.Fatal("Expected 'read' to be TypeList")
	}
	if r.Schema["create"].Type != schema.TypeList {
		t.Fatal("Expected 'create' to be TypeList")
	}
	if r.Schema["update"].Type != schema.TypeList {
		t.Fatal("Expected 'update' to be TypeList")
	}
	if r.Schema["delete"].Type != schema.TypeList {
		t.Fatal("Expected 'delete' to be TypeList")
	}
	if r.Schema["target_state"].Type != schema.TypeList {
		t.Fatal("Expected 'target_state' to be TypeList")
	}
	if r.Schema["working_dir"].Type != schema.TypeString {
		t.Fatal("Expected 'working_dir' to be TypeString")
	}
	if r.Schema["resource"].Type != schema.TypeString {
		t.Fatal("Expected 'resource' to be TypeString")
	}
}

func TestAccResourceScript_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC set")
	}

	// Create temporary test files
	tmpDir, err := ioutil.TempDir("", "terraform-provider-script-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create script files for testing
	createScript := filepath.Join(tmpDir, "create.sh")
	readScript := filepath.Join(tmpDir, "read.sh")
	updateScript := filepath.Join(tmpDir, "update.sh")
	deleteScript := filepath.Join(tmpDir, "delete.sh")
	targetStateScript := filepath.Join(tmpDir, "target_state.sh")

	// Write test scripts
	createScriptContent := `#!/bin/sh
echo '{"id":"test-resource-id","resource":"initial-state"}'
`
	readScriptContent := `#!/bin/sh
echo '{"id":"test-resource-id","resource":"initial-state"}'
`
	updateScriptContent := `#!/bin/sh
echo '{"id":"test-resource-id","resource":"updated-state"}'
`
	deleteScriptContent := `#!/bin/sh
exit 0
`
	targetStateScriptContent := `#!/bin/sh
echo '{"data":"test-data"}'
`

	if err := ioutil.WriteFile(createScript, []byte(createScriptContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(readScript, []byte(readScriptContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(updateScript, []byte(updateScriptContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(deleteScript, []byte(deleteScriptContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(targetStateScript, []byte(targetStateScriptContent), 0755); err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"script": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScriptConfig(tmpDir, createScript, readScript, updateScript, deleteScript, targetStateScript),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("script.test", "resource", "initial-state"),
				),
			},
		},
	})
}

func testAccResourceScriptConfig(workingDir, createScript, readScript, updateScript, deleteScript, targetStateScript string) string {
	return fmt.Sprintf(`
resource "script" "test" {
  working_dir = "%s"
  
  create = ["sh", "%s"]
  read = ["sh", "%s", "##ID##"]
  update = ["sh", "%s", "##ID##", "##CS##"]
  delete = ["sh", "%s", "##ID##"]
  target_state = ["sh", "%s"]
}
`, 
		workingDir,
		createScript,
		readScript,
		updateScript,
		deleteScript,
		targetStateScript,
	)
} 
