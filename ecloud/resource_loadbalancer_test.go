package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoadBalancer_basic(t *testing.T) {
	lbName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_loadbalancer.test-lb"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLoadBalancerConfig_basic(ANS_TEST_VPC_REGION_ID, lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", lbName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetLoadBalancer(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.LoadBalancerNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckLoadBalancerDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_loadbalancer" {
			continue
		}

		_, err := service.GetLoadBalancer(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Loadbalancer with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.LoadBalancerNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceLoadBalancerConfig_basic(regionID string, lbName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_loadbalancer_spec" "medium-lb" {
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
	name = "%s"
	load_balancer_spec_id = data.ecloud_loadbalancer_spec.medium-lb.id
	network_id = ecloud_network.test-network.id
}
`, regionID, lbName)
}
