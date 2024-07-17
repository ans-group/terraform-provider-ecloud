package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNService_basic(t *testing.T) {
	vpnServiceName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPNServiceConfig_basic(vpnServiceName)
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

func testAccDataSourceVPNServiceConfig_basic(vpnServiceName string) string {
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

resource "ecloud_vpn_service" "test-vpnservice" {
	router_id = ecloud_router.test-router.id
	name = "%[1]s"
}

data "ecloud_vpn_service" "test-vpnservice" {
	name = "%[1]s"
}
`, vpnServiceName)
}
