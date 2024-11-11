package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNGatewayUser_basic(t *testing.T) {
	vpnGatewayUserName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNGatewayUserConfig_basic(vpnGatewayUserName)
	resourceName := "data.ecloud_vpn_gateway_user.test-vpngatewayuser"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayUserName),
				),
			},
		},
	})
}

func testAccDataSourceVPNGatewayUserConfig_basic(vpnGatewayUserName string) string {
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

data "ecloud_vpn_gateway_specification" "test-spec" {
	name = "Small"
}

resource "ecloud_vpn_gateway" "test-vpngateway" {
	router_id = ecloud_router.test-router.id
	name = "tftest-vpngateway"
	specification_id = data.ecloud_vpn_gateway_specification.test-spec.id
}

resource "ecloud_vpn_gateway_user" "test-vpngatewayuser" {
	vpn_gateway_id = ecloud_vpn_gateway.test-vpngateway.id
	name = "%[1]s"
	username = "tftest-user"
	password = "password123!"
}

data "ecloud_vpn_gateway_user" "test-vpngatewayuser" {
	name = "%[1]s"
}
`, vpnGatewayUserName)
}
