package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVPNSession_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")

	resourceName := "ecloud_vpn_service.test-vpnservice"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNSessionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNSessionConfig_basic(UKF_TEST_VPC_REGION_ID, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNSessionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccCheckVPNSessionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN session ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPNSession(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPNSessionNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPNSessionDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpn_service" {
			continue
		}

		_, err := service.GetVPNSession(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN session with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPNSessionNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVPNSessionConfig_basic(regionID string, vpcName string) string {
	return fmt.Sprintf(`
	resource "ecloud_vpc" "test-vpc" {
		region_id = "%s"
		name = "%s"
	}
`, regionID, vpcName)
}
