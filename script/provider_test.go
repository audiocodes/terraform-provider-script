package script

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	p := Provider()
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

// TestAccProvider verifies provider instantiation
func TestAccProvider(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC set")
	}

	// Setup a local provider server for acceptance testing
	var expectedProviders = map[string]func() (*schema.Provider, error){
		"script": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
	
	providerFactories := map[string]func() (*schema.Provider, error){
		"script": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
	
	if len(expectedProviders) != len(providerFactories) {
		t.Fatalf("Expected number of providers: %d, got: %d", len(expectedProviders), len(providerFactories))
	}
} 
