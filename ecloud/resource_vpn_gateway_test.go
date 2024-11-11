package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVPNGateway_basic(t *testing.T) {
	vpnGatewayName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_vpn_gateway.test-vpngateway"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNGatewayConfig_basic(vpnGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayName),
					resource.TestCheckResourceAttrSet(resourceName, "fqdn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPNGateway_update(t *testing.T) {
	vpnGatewayName := acctest.RandomWithPrefix("tftest")
	vpnGatewayNameUpdated := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_vpn_gateway.test-vpngateway"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNGatewayConfig_basic(vpnGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayName),
				),
			},
			{
				Config: testAccResourceVPNGatewayConfig_basic(vpnGatewayNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayNameUpdated),
				),
			},
		},
	})
}

func testAccCheckVPNGatewayExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN gateway ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPNGateway(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPNGatewayNotFoundError); ok {
				return fmt.Errorf("VPN gateway not found")
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPNGatewayDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpn_gateway" {
			continue
		}

		_, err := service.GetVPNGateway(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN gateway with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPNGatewayNotFoundError); ok {
			continue
		}

		return err
	}

	return nil
}

func testAccResourceVPNGatewayConfig_basic(gatewayName string) string {
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
	name = "%[1]s"
	specification_id = data.ecloud_vpn_gateway_specification.test-spec.id
}
`, gatewayName)
}
