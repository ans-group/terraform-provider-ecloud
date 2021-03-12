package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccInstance_basic(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_instance.test-instance"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstanceConfig_basic(UKF_TEST_VPC_REGION_ID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetInstance(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.InstanceNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_instance" {
			continue
		}

		_, err := service.GetInstance(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Instance with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.InstanceNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceInstanceConfig_basic(regionID string, instanceName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "%s"
	image_id = "img-abcdef"
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}
`, regionID, instanceName)
}
