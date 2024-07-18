package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFloatingIP_basic(t *testing.T) {
	fipName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceFloatingIPConfig_basic(fipName)
	resourceName := "data.ecloud_floatingip.test-fip"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
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

func testAccDataSourceFloatingIPConfig_basic(fipName string) string {
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

resource "ecloud_floatingip" "test-fip" {
   vpc_id = ecloud_vpc.test-vpc.id
   availability_zone_id = data.ecloud_availability_zone.test-az.id
   name = "%[1]s"
 }

data "ecloud_floatingip" "test-fip" {
    name = ecloud_floatingip.test-fip.name
}
`, fipName)
}
