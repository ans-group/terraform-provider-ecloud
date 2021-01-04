package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccVPC_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")

	resourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVPCConfig_basic(UKF_TEST_VPC_REGION_ID, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccCheckVPCExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVPC(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPCNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVPCDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_vpc" {
			continue
		}

		_, err := service.GetVPC(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VPCNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVPCConfig_basic(regionID string, vpcName string) string {
	return fmt.Sprintf(`
	resource "ecloud_vpc" "test-vpc" {
		region_id = "%s"
		name = "%s"
	}
`, regionID, vpcName)
}
