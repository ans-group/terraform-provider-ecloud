package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNSession_basic(t *testing.T) {
	vpnSessionName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNSessionConfig_basic(vpnSessionName)
	resourceName := "data.ecloud_vpn_session.test-vpnsession"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", vpnSessionName),
				),
			},
		},
	})
}

func testAccDataSourceVPNSessionConfig_basic(vpnSessionName string) string {
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
	name = "tftest-router"
	availability_zone_id = data.ecloud_availability_zone.test-az.id
}

resource "ecloud_vpn_service" "test-vpnservice" {
	router_id = ecloud_router.test-router.id
	name = "tftest-vpnservice"
}

resource "ecloud_vpn_endpoint" "test-vpnendpoint" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	name = "tftest-vpnendpoint"
}

resource "ecloud_vpn_session" "test-vpnsession" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	vpn_endpoint_id = ecloud_vpn_endpoint.test-vpnendpoint.id
	remote_ip = "1.2.3.4"
	name = "%[1]s"
}

data "ecloud_vpn_session" "test-vpnsession" {
	name = "%[1]s"
}
`, vpnSessionName)
}
