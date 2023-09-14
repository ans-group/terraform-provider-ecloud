package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNATOverloadRule_basic(t *testing.T) {
	ruleName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_natoverloadrule.test-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNATOverloadRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNATOverloadRuleConfig_basic(ANS_TEST_VPC_REGION_ID, ruleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATOverloadRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", ruleName),
				),
			},
		},
	})
}

func testAccCheckNATOverloadRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NATOverloadRule ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetNATOverloadRule(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NATOverloadRuleNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckNATOverloadRuleDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_network" {
			continue
		}

		_, err := service.GetNATOverloadRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("NATOverloadRule with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.NATOverloadRuleNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceNATOverloadRuleConfig_basic(regionID string, ruleName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_floatingip" "test-fip" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	resource_id = ecloud_router.test-router.id
}
  
resource "ecloud_natoverloadrule" "test-rule" {
	name = "%[2]s"
	network_id = ecloud_network.test-network.id
	subnet = "10.0.0.0/24"
	floating_ip_id = ecloud_floatingip.test-fip.id
	action = "allow"
}
`, regionID, ruleName)
}
