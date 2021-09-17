package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNSession_basic(t *testing.T) {
	vpnSessionName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNSessionConfig_basic(UKF_TEST_VPC_REGION_ID, vpnSessionName)
	resourceName := "data.ecloud_vpn_session.test-vpnsession"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func testAccDataSourceVPNSessionConfig_basic(regionID string, vpnSessionName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
	availability_zone_id = data.ecloud_availability_zone.test-az.id
}

resource "ecloud_vpn_service" "test-vpnservice" {
	router_id = ecloud_router.test-router.id
	name = "test-vpnservice"
}

resource "ecloud_vpn_endpoint" "test-vpnendpoint" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	name = "test-vpnendpoint"
}

resource "ecloud_vpn_session" "test-vpnsession" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	vpn_endpoint_id = ecloud_vpn_endpoint.test-vpnendpoint.id
	remote_ip = "1.2.3.4"
	name = "%[2]s"
}

data "ecloud_vpn_session" "test-vpnsession" {
	name = "%[2]s"
}
`, regionID, vpnSessionName)
}
