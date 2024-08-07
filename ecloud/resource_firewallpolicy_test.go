package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewallPolicy_basic(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_firewallpolicy.test-fwp"
	routerResourceName := "ecloud_router.test-router"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFirewallPolicyConfig_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttrPair(routerResourceName, "id", resourceName, "router_id"),
					resource.TestCheckResourceAttr(resourceName, "sequence", "0"),
				),
			},
		},
	})
}

func testAccCheckFirewallPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall policy ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetFirewallPolicy(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.FirewallPolicyNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckFirewallPolicyDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_firewallpolicy" {
			continue
		}

		_, err := service.GetFirewallPolicy(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Firewall policy with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.FirewallPolicyNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceFirewallPolicyConfig_basic(policyName string) string {
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

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_firewallpolicy" "test-fwp" {
	router_id = ecloud_router.test-router.id
	name = "%[1]s"
	sequence = 0
}
`, policyName)
}
