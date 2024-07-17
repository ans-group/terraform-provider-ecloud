package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoadBalancerVip_basic(t *testing.T) {
	VIPName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_loadbalancer_vip.lb-vip"
	lbResourceName := "ecloud_loadbalancer.test-lb"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLoadBalancerVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLoadBalancerVipConfig_basic(VIPName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerVipExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", VIPName),
					resource.TestCheckResourceAttrPair(lbResourceName, "id", resourceName, "load_balancer_id"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerVipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer vip ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVIP(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VIPNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckLoadBalancerVipDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_loadbalancer_vip" {
			continue
		}

		_, err := service.GetVIP(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Loadbalancer vip with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VIPNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceLoadBalancerVipConfig_basic(VIPName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
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
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
	subnet = "10.0.1.0/24"
}

resource "ecloud_loadbalancer" "test-lb" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-lb"
	load_balancer_spec_id = data.ecloud_loadbalancer_spec.medium-lb.id
	network_id = ecloud_network.test-network.id
}

resource "ecloud_loadbalancer_vip" "lb-vip" {
	name = "%[1]s"
	load_balancer_id = data.ecloud_loadbalancer.test-lb.id
}
`, VIPName)
}
