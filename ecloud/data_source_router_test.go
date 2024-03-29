package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRouter_basic(t *testing.T) {
	routerName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceRouterConfig_basic(ANS_TEST_VPC_REGION_ID, routerName)
	resourceName := "data.ecloud_router.test-router"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", routerName),
				),
			},
		},
	})
}

func testAccDataSourceRouterConfig_basic(regionID string, routerName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name      = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name   = "%s"
}

data "ecloud_router" "test-router" {
    name = ecloud_router.test-router.name
}
`, regionID, routerName)
}
