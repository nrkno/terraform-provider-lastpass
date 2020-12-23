package lastpass

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"lastpass": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LASTPASS_USER"); v == "" {
		t.Fatal("LASTPASS_USER must be set for acceptance tests")
	}
	if v := os.Getenv("LASTPASS_PASSWORD"); v == "" {
		t.Fatal("LASTPASS_PASSWORD must be set for acceptance tests")
	}
}
