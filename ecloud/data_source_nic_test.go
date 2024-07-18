package ecloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNic_basic(t *testing.T) {
	config := testAccDataSourceNicConfig_basic()
	resourceName := "data.ecloud_nic.test-nic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "ip_address", regexp.MustCompile(`^(10\.0\.0\.)(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)),
				),
			},
		},
	})
}

func testAccDataSourceNicConfig_basic() string {
	return `
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
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
	name = "tftest-instance"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 40
	ram_capacity = 1024
	vcpu_cores = 1
}

data "ecloud_nic" "test-nic" {
    nic_id = ecloud_instance.test-instance.nic_id
}
`
}
