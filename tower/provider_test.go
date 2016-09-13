package tower

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"tower": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if e := os.Getenv("TOWER_ENDPOINT"); e == "" {
		t.Fatal("TOWER_ENDPOINT must be set for acceptance tests")
	}
	if u := os.Getenv("TOWER_USERNAME"); u == "" {
		t.Fatal("TOWER_USERNAME must be set for acceptance tests")
	}
	if p := os.Getenv("TOWER_PASSWORD"); p == "" {
		t.Fatal("TOWER_PASSWORD must be set for acceptance tests")
	}
}
