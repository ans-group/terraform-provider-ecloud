package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkPolicy_basic(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_networkpolicy.test-np"
	routerResourceName := "ecloud_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkPolicyConfig_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttrPair(routerResourceName, "id", resourceName, "network_id"),
					resource.TestCheckResourceAttr(resourceName, "catchall_rule_action", "REJECT"),
				),
			},
		},
	})
}

func testAccCheckNetworkPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network policy ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetNetworkPolicy(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NetworkPolicyNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckNetworkPolicyDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_networkpolicy" {
			continue
		}

		_, err := service.GetNetworkPolicy(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network policy with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.NetworkPolicyNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceNetworkPolicyConfig_basic(policyName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
	advanced_networking = true
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_networkpolicy" "test-np" {
	network_id = ecloud_network.test-network.id
	name = "%[1]s"
	catchall_rule_action = "REJECT"
}
`, policyName)
}
