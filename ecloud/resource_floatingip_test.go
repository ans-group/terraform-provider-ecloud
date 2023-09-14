package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFloatingIP_basic(t *testing.T) {
	fipName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_floatingip.test-fip"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFloatingIPConfig_basic(ANS_TEST_VPC_REGION_ID, fipName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFloatingIPExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fipName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckFloatingIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No floating ip ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetFloatingIP(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.FloatingIPNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckFloatingIPDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_floatingip" {
			continue
		}

		_, err := service.GetFloatingIP(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Floating ip with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.FloatingIPNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceFloatingIPConfig_basic(regionID string, fipName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_floatingip" "test-fip" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%s"
}
`, regionID, fipName)
}
