package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceFloatingIP_basic(t *testing.T) {
	fipName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceFloatingIPConfig_basic(UKF_TEST_VPC_REGION_ID, fipName)
	resourceName := "data.ecloud_floatingip.test-fip"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fipName),
				),
			},
		},
	})
}

func testAccDataSourceFloatingIPConfig_basic(regionID string, fipName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_floatingip" "test-fip" {
   vpc_id = ecloud_vpc.test-vpc.id
   availability_zone_id = data.ecloud_availability_zone.test-az.id
   name = "%[2]s"
 }

data "ecloud_floatingip" "test-fip" {
    name = ecloud_floatingip.test-fip.name
}
`, regionID, fipName)
}
