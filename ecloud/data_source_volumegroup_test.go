package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVolumeGroup_basic(t *testing.T) {
	volumeGroupName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVolumeGroupConfig_basic(ANS_TEST_VPC_REGION_ID, volumeGroupName)
	resourceName := "data.ecloud_volumegroup.test-volumegroup"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", volumeGroupName),
				),
			},
		},
	})
}

func testAccDataSourceVolumeGroupConfig_basic(regionID string, volumeGroupName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_volumegroup" "test-volumegroup" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%[2]s"
}

data "ecloud_volumegroup" "test-volumegroup" {
    name = ecloud_volumegroup.test-volumegroup.name
}
`, regionID, volumeGroupName)
}
