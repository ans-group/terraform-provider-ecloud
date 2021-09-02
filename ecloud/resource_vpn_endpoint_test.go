package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccVPNEndpoint_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")

	resourceName := "ecloud_vpn_service.test-vpnservice"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNEndpointConfig_basic(UKF_TEST_VPC_REGION_ID, vpcName),
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

func testAccResourceVPNEndpointConfig_basic(regionID string, vpcName string) string {
	return fmt.Sprintf(`
	resource "ecloud_vpc" "test-vpc" {
		region_id = "%s"
		name = "%s"
	}
`, regionID, vpcName)
}
