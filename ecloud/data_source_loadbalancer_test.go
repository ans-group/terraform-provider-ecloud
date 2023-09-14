package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLoadBalancer_basic(t *testing.T) {
	lbName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceLoadBalancerConfig_basic(ANS_TEST_VPC_REGION_ID, lbName)
	resourceName := "data.ecloud_loadbalancer.test-lb"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", lbName),
				),
			},
		},
	})
}

func testAccDataSourceLoadBalancerConfig_basic(regionID string, lbName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_loadbalancer_spec" "test-lb-medium" {
	name = "Medium"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
	subnet = "10.0.1.0/24"
}

resource "ecloud_loadbalancer" "test-lb" {
   vpc_id = ecloud_vpc.test-vpc.id
   availability_zone_id = data.ecloud_availability_zone.test-az.id
   name = "%[2]s"
   load_balancer_spec_id = data.ecloud_loadbalancer_spec.test-lb-medium.id
   network_id = ecloud_network.test-network.id
}

data "ecloud_loadbalancer" "test-lb" {
    name = ecloud_loadbalancer.test-lb.name
}
`, regionID, lbName)
}
