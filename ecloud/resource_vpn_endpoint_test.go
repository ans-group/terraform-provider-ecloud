package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVPNEndpoint_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")

	resourceName := "ecloud_vpn_service.test-vpnservice"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVPNEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNEndpointConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccCheckVPNEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN endpoint ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPNEndpoint(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPNEndpointNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPNEndpointDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpn_service" {
			continue
		}

		_, err := service.GetVPNEndpoint(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN endpoint with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPNEndpointNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVPNEndpointConfig_basic(vpcName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "%s"
}
`, vpcName)
}
