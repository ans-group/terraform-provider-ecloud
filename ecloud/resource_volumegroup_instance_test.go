package ecloud

import (
	"errors"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolumeGroupInstance_basic(t *testing.T) {
	resourceName := "ecloud_volumegroup_instance.test-volumegroup-instance"
	instanceResourceName := "ecloud_vpc.test-instance"
	volumeGroupResourceName := "ecloud_vpc.test-instance"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVolumeGroupInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceVolumeGroupInstanceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(instanceResourceName, "id", resourceName, "instance_id"),
					resource.TestCheckResourceAttrPair(volumeGroupResourceName, "id", resourceName, "volume_group_id"),
				),
			},
		},
	})
}

func testAccCheckVolumeGroupInstanceDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_volumegroup_instance" {
			continue
		}

		instanceID := rs.Primary.Attributes["instance_id"]

		instance, err := service.GetInstance(instanceID)
		if err != nil {
			return errors.New("Failed to retrieve instance")
		}

		if instance.VolumeGroupID != "" {
			return errors.New("Volume group still attached to instance")
		}

		return nil
	}

	return nil
}

func testAccResourceVolumeGroupInstanceConfig_basic() string {
	return `
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_image" "almalinux9" {
	name = "AlmaLinux 9"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	image_id = data.ecloud_image.almalinux9.id
	volume_capacity = 40
	ram_capacity = 1024
	vcpu {
		sockets          = 1
		cores_per_socket = 2
	}
}

resource "ecloud_volumegroup" "test-volumegroup" {
    vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
}

resource "ecloud_volume" "test-volume" {
	vpc_id               = ecloud_vpc.test-vpc.id
	capacity             = 10
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	volume_group_id      = ecloud_volumegroup.test-volumegroup.id
}

resource "ecloud_volumegroup_instance" "test-volume-group-instance" {
	volume_group_id = ecloud_volumegroup.test-volumegroup.id
	instance_id     = ecloud_instance.test-instance.id
}
`
}
