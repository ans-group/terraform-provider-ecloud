package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVolume_basic(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVolumeConfig_basic(ANS_TEST_VPC_REGION_ID, volumeName)
	resourceName := "data.ecloud_volume.test-volume"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", volumeName),
				),
			},
		},
	})
}

func testAccDataSourceVolumeConfig_basic(regionID string, volumeName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_volume" "test-volume" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%[2]s"
	capacity = 1
	iops = 300
}

data "ecloud_volume" "test-volume" {
    name = ecloud_volume.test-volume.name
}
`, regionID, volumeName)
}
