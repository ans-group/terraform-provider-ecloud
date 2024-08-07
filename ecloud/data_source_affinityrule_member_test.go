package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAffinityRuleMember_basic(t *testing.T) {
	config := testAccDataSourceAffinityRuleMemberConfig_basic()
	armResourceName := "data.affinityrule_member.test-arm"
	arResourceName := "ecloud_affinityrule.test-ar.id"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(armResourceName, "affinity_rule_id", arResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceAffinityRuleMemberConfig_basic() string {
	return `
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

data "ecloud_image" "centos7" {
	name = "CentOS 7"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
	subnet = "10.0.1.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "tftest-instance"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

resource "ecloud_affinityrule" "test-ar" {
   vpc_id = ecloud_vpc.test-vpc.id
   availability_zone_id = data.ecloud_availability_zone.test-az.id
   name = "tftest-ar"
   type = "anti-affinity"
}

resource "ecloud_affinityrule_member" "test-arm" {
	affinity_rule_id = ecloud_affinityrule.test-ar.id
	instance_id = ecloud_instance.test-instance.id
}

data "ecloud_affinityrule_member" "test-arm" {
    instance_id = ecloud_instance.test-instance.id
	affinity_rule_id = ecloud_affinityrule.test-ar.id
}
`
}
