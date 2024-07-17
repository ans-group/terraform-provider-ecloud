package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceInstance_basic(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceInstanceConfig_basic(instanceName)
	resourceName := "data.ecloud_instance.test-instance"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
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

func testAccDataSourceInstanceConfig_basic(instanceName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
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
	name = "test-router"
	availability_zone_id = data.ecloud_availability_zone.test-az.id
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "%[1]s"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 40
	ram_capacity = 1024
	vcpu_cores = 1
}

data "ecloud_instance" "test-instance" {
    name = ecloud_instance.test-instance.name
}
`, instanceName)
}
