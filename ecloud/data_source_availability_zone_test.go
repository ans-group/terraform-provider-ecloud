package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAvailabilityZone_basic(t *testing.T) {
	azName := "Manchester South"
	config := testAccDataSourceAvailabilityZoneConfig_basic(azName)
	resourceName := "data.ecloud_availability_zone.test-az"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", azName),
				),
			},
		},
	})
}

func testAccDataSourceAvailabilityZoneConfig_basic(azName string) string {
	return fmt.Sprintf(`
data "ecloud_availability_zone" "test-az" {
    name = "%s"
}
`, azName)
}
