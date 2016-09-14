package tower

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccArtifact_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInventory_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArtifactState("name", "tf-provider-test19"),
					testAccCheckArtifactState("organization_id", "1"),
				),
			},
		},
	})
}

func TestAccArtifact_metadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInventory_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckArtifactState("name", "tf-provider-test20"),
					testAccCheckArtifactState("organization_id", "1"),
					testAccCheckArtifactState("description", "foobar"),
				),
			},
		},
	})
}

func testAccCheckArtifactState(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["tower_inventory.foobar"]
		if !ok {
			return fmt.Errorf("Not found: %s", "tower_inventory.foobar")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		p := rs.Primary
		if p.Attributes[key] != value {
			return fmt.Errorf(
				"%s != %s (actual: %s)", key, value, p.Attributes[key])
		}

		return nil
	}
}

const testAccInventory_basic = `
resource "tower_inventory" "foobar" {
	name = "tf-provider-test19"
	organization_id = "1"
}`

const testAccInventory_description = `
resource "tower_inventory" "foobar" {
	name = "tf-provider-test20"
	organization_id = "1"
	description = "foobar"
}`
