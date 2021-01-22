package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNetwork_basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("tftest")
	subnet := "10.0.0.0/24"
	config := testAccDataSourceNetworkConfig_basic(UKF_TEST_VPC_REGION_ID, networkName, subnet)
	resourceName := "data.ecloud_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", networkName),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnet),
				),
			},
		},
	})
}

func testAccDataSourceNetworkConfig_basic(regionID string, networkName string, subnet string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "%[2]s"
	subnet = "%[3]s"
}

data "ecloud_network" "test-network" {
    name = ecloud_network.test-network.name
}
`, regionID, networkName, subnet)
}
