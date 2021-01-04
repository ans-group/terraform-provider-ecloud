package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNetwork_basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceNetworkConfig_basic(UKF_TEST_VPC_REGION_ID, networkName)
	resourceName := "data.ecloud_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", networkName),
				),
			},
		},
	})
}

func testAccDataSourceNetworkConfig_basic(regionID string, networkName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "%s"
}

data "ecloud_network" "test-network" {
    name = ecloud_network.test-network.name
}
`, regionID, networkName)
}
