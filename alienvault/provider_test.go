package alienvault

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {

	required := []string{
		"ALIENVAULT_FQDN",
		"ALIENVAULT_USERNAME",
		"ALIENVAULT_PASSWORD",
	}

	for _, req := range required {
		if v := os.Getenv(req); v == "" {
			t.Fatalf("%s must be set for acceptance tests", req)
		}
	}
}

func init() {
	_ = os.Setenv("ALIENVAULT_SKIP_TLS_VERIFY", "1")
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"alienvault": testAccProvider,
	}
}
