package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRegion_basic(t *testing.T) {
	region := "London"
	config := testAccDataSourceRegionConfig_basic(region)
	resourceName := "data.ecloud_region.lon-region"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", region),
				),
			},
		},
	})
}

func testAccDataSourceRegionConfig_basic(region string) string {
	return fmt.Sprintf(`
data "ecloud_region" "lon-region" {
    name = "%s"
}
`, region)
}
