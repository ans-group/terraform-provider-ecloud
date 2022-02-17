package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccLoadBalancerNetwork_basic(t *testing.T) {
	lbName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_loadbalancer_network.lb-network"
	lbResourceName := "ecloud_loadbalancer.test-lb"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLoadBalancerNetworkConfig_basic(UKF_TEST_VPC_REGION_ID, lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", lbName),
					resource.TestCheckResourceAttrPair(lbResourceName, "id", resourceName, "load_balancer_id"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer network ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetLoadBalancerNetwork(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.LoadBalancerNetworkNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckLoadBalancerNetworkDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_loadbalancer_network" {
			continue
		}

		_, err := service.GetLoadBalancerNetwork(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Loadbalancer network with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.LoadBalancerNetworkNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceLoadBalancerNetworkConfig_basic(regionID string, lbName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_loadbalancer_spec" "medium-lb" {
	name = "Medium
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
	name = "test-lb"
	load_balancer_spec_id = data.ecloud_loadbalancer_spec.medium-lb.id
}

resource "ecloud_loadbalancer_network" "lb-network" {
	network_id= ecloud_network.test-network.id
	name = "%[2]s"
	load_balancer_id = data.ecloud_loadbalancer.test-lb.id
}
`, regionID, lbName)
}
