package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceRouter_basic(t *testing.T) {
	routerName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceRouterConfig_basic(UKF_TEST_VPC_REGION_ID, routerName)
	resourceName := "data.ecloud_router.test-router"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name   = "%s"
}

data "ecloud_router" "test-router" {
    name = ecloud_router.test-router.name
}
`, regionID, routerName)
}
