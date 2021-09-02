package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVPNEndpoint_basic(t *testing.T) {
	vpnEndpointName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNEndpointConfig_basic(UKF_TEST_VPC_REGION_ID, vpnEndpointName)
	resourceName := "data.ecloud_vpn_endpoint.test-vpnendpoint"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func testAccDataSourceVPNEndpointConfig_basic(regionID string, vpnEndpointName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_vpn_service" "test-vpnservice" {
	router_id = ecloud_router.test-router.id
	name = "test-vpnservice"
}

resource "ecloud_vpn_endpoint" "test-vpnendpoint" {
	vpn_service_id = ecloud_vpn_service.test-vpnservice.id
	name = "%[2]s"
}

data "ecloud_vpn_endpoint" "test-vpnendpoint" {
	name = "%[2]s"
}
`, regionID, vpnEndpointName)
}
