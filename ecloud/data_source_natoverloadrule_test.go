package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNATOverloadRule_basic(t *testing.T) {
	ruleName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceNATOverloadRuleConfig_basic(ANS_TEST_VPC_REGION_ID, ruleName)
	resourceName := "data.ecloud_natoverloadrule.test-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", ruleName),
					resource.TestCheckResourceAttr(resourceName, "subnet", "10.0.0.0/24"),
				),
			},
		},
	})
}

func testAccDataSourceNATOverloadRuleConfig_basic(regionID string, ruleName string) string {
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

data "ecloud_natoverloadrule" "test-rule" {
    name = ecloud_natoverloadrule.test-rule.name
}
`, regionID, ruleName)
}
