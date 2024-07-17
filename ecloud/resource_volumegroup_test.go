package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolumeGroup_basic(t *testing.T) {
	volumeGroupName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_volumegroup.test-volumegroup"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVolumeGroupConfig_basic(volumeGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", volumeGroupName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckVolumeGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volumegroup ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetVolumeGroup(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VolumeGroupNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckVolumeGroupDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_volumegroup" {
			continue
		}

		_, err := service.GetVolumeGroup(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Volumegroup with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.VolumeGroupNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceVolumeGroupConfig_basic(volumeGroupName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_volumegroup" "test-volumegroup" {
    vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
    name = "%[1]s"
}
`, volumeGroupName)
}
