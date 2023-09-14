package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNService_basic(t *testing.T) {
	vpnServiceName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNServiceConfig_basic(ANS_TEST_VPC_REGION_ID, vpnServiceName)
	resourceName := "data.ecloud_vpn_service.test-vpnservice"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", vpnServiceName),
				),
			},
		},
	})
}

func testAccDataSourceVPNServiceConfig_basic(regionID string, vpnServiceName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
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
	name = "%[2]s"
}

data "ecloud_vpn_service" "test-vpnservice" {
	name = "%[2]s"
}
`, regionID, vpnServiceName)
}
