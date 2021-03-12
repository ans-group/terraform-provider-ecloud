package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceInstance_basic(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceInstanceConfig_basic(UKF_TEST_VPC_REGION_ID, instanceName)
	resourceName := "data.ecloud_instance.test-instance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
				),
			},
		},
	})
}

func testAccDataSourceInstanceConfig_basic(regionID string, instanceName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "%[2]s"
	image_id = "img-abcdef12"
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

data "ecloud_instance" "test-instance" {
    name = ecloud_instance.test-instance.name
}
`, regionID, instanceName)
}
