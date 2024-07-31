package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceECloudInstanceCreds(t *testing.T) {
	credName := "root"
	config := testAccDataSourceECloudInstanceCreds_basic(credName)
	resourceName := "data.ecloud_instance_credential.root_creds"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", credName),
				),
			},
		},
	})
}

func testAccDataSourceECloudInstanceCreds_basic(credName string) string {
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

resource "ecloud_router" "test-router-1" {
  vpc_id               = ecloud_vpc.test-vpc.id
  availability_zone_id = data.ecloud_availability_zone.test-az.id
  name                 = "test-router"
}

resource "ecloud_network" "network-1" {
	router_id = ecloud_router.test-router-1.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "instance-1" {
  vcpu {
    sockets          = 1
    cores_per_socket = 2
  }
  ram_capacity    = 2048
  vpc_id          = ecloud_vpc.test-vpc.id
  name            = "instance test"
  image_id        = "img-19cb94e5"
  volume_capacity = 40
  volume_iops     = 600
  network_id      = ecloud_network.network-1.id
  backup_enabled  = false
  encrypted       = false
}

data "ecloud_instance_credential" "root_creds" {
  instance_id = ecloud_instance.instance-1.id
  name        = "%[1]s"
}
`, credName)
}
