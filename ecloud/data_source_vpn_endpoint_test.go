package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNEndpoint_basic(t *testing.T) {
	vpnEndpointName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNEndpointConfig_basic(vpnEndpointName)
	resourceName := "data.ecloud_vpn_endpoint.test-vpnendpoint"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", vpnEndpointName),
				),
			},
		},
	})
}

func testAccDataSourceVPNEndpointConfig_basic(vpnEndpointName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
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

resource "ecloud_vpn_service" "test-vpnservice" {
	router_id = ecloud_router.test-router.id
	name = "test-vpnservice"
}

resource "ecloud_vpn_endpoint" "test-vpnendpoint" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	name = "%[1]s"
}

data "ecloud_vpn_endpoint" "test-vpnendpoint" {
	name = "%[1]s"
}
`, vpnEndpointName)
}
