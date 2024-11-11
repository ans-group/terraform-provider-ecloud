package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVPNGatewayUser_basic(t *testing.T) {
	vpnGatewayUserName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_vpn_gateway_user.test-vpngatewayuser"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPNGatewayUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNGatewayUserConfig_basic(vpnGatewayUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayUserName),
					resource.TestCheckResourceAttr(resourceName, "username", "tftest-user"),
					resource.TestCheckResourceAttr(resourceName, "password", "password123!"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func TestAccVPNGatewayUser_update(t *testing.T) {
	vpnGatewayUserName := acctest.RandomWithPrefix("tftest")
	vpnGatewayUserNameUpdated := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_vpn_gateway_user.test-vpngatewayuser"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPNGatewayUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNGatewayUserConfig_basic(vpnGatewayUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayUserName),
					resource.TestCheckResourceAttr(resourceName, "password", "password123!"),
				),
			},
			{
				Config: testAccResourceVPNGatewayUserConfig_update(vpnGatewayUserNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpnGatewayUserNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "password", "newpassword123!"),
				),
			},
		},
	})
}

func testAccCheckVPNGatewayUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN gateway user ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPNGatewayUser(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPNGatewayUserNotFoundError); ok {
				return fmt.Errorf("VPN gateway user not found")
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPNGatewayUserDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpn_gateway_user" {
			continue
		}

		_, err := service.GetVPNGatewayUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN gateway user with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPNGatewayUserNotFoundError); ok {
			continue
		}

		return err
	}

	return nil
}

func testAccResourceVPNGatewayUserConfig_basic(userName string) string {
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
`, userName)
}

func testAccResourceVPNGatewayUserConfig_update(userName string) string {
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
	password = "newpassword123!"
}
`, userName)
}
