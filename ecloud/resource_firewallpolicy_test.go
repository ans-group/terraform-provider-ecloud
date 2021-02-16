package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccFirewallPolicy_basic(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_firewallpolicy.test-fwp"
	routerResourceName := "ecloud_router.test-router"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFirewallPolicyConfig_basic(UKF_TEST_VPC_REGION_ID, policyName),
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

func testAccResourceFirewallPolicyConfig_basic(regionID string, policyName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_firewallpolicy" "test-fwp" {
	router_id = ecloud_router.test-router.id
	name = "%[2]s"
	sequence = 0
}
`, regionID, policyName)
}
