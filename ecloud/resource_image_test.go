package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImage_basic(t *testing.T) {
	imageName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_image.test-image"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckimageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceImageConfig_basic(UKF_TEST_VPC_REGION_ID, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckimageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", imageName),
				),
			},
		},
	})
}

func testAccCheckimageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetImage(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.ImageNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckimageDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_image" {
			continue
		}

		_, err := service.GetImage(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Image with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.ImageNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceImageConfig_basic(regionID string, imageName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

data "ecloud_image" "centos7" {
	name = "CentOS 7"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "test-instance"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

resource "ecloud_image" "test-image" {
	instance_id = ecloud_instance.test-instance.id
	name = "%s"
}
`, regionID, imageName)
}
