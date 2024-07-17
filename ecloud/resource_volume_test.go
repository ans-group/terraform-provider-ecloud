package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolume_basic(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_volume.test-volume"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVolumeConfig_basic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", volumeName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVolume(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VolumeNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVolumeDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_volume" {
			continue
		}

		_, err := service.GetVolume(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Volume with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VolumeNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVolumeConfig_basic(volumeName string) string {
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

resource "ecloud_volume" "test-volume" {
    vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
    capacity = 1
    name = "%[1]s"
    iops = 300
}
`, volumeName)
}
