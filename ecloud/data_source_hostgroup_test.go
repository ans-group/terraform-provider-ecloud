package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHostGroup_basic(t *testing.T) {
	hostGroupName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceHostGroupConfig_basic(hostGroupName)
	resourceName := "data.ecloud_hostgroup.test-hostgroup"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", hostGroupName),
				),
			},
		},
	})
}

func testAccDataSourceHostGroupConfig_basic(hostGroupName string) string {
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

data "ecloud_hostspec" "test-hostspec" {
	name = "DUAL-4208--64GB"
}

resource "ecloud_hostgroup" "test-hostgroup" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	host_spec_id = data.ecloud_hostspec.test-hostspec.id
	name = "%[1]s"
    windows_enabled = false
}

data "ecloud_hostgroup" "test-hostgroup" {
    name = ecloud_hostgroup.test-hostgroup.name
}
`, hostGroupName)
}
