package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccVPNService_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")

	resourceName := "ecloud_vpn_service.test-vpnservice"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPNServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPNServiceConfig_basic(UKF_TEST_VPC_REGION_ID, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccCheckVPNServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN service ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPNService(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPNServiceNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPNServiceDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpn_service" {
			continue
		}

		_, err := service.GetVPNService(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN service with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPNServiceNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVPNServiceConfig_basic(regionID string, vpcName string) string {
	return fmt.Sprintf(`
	resource "ecloud_vpc" "test-vpc" {
		region_id = "%s"
		name = "%s"
	}
`, regionID, vpcName)
}
