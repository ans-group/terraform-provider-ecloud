package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHost_basic(t *testing.T) {
	hostName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceHostConfig_basic(ANS_TEST_VPC_REGION_ID, hostName)
	resourceName := "data.ecloud_host.test-host"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", hostName),
				),
			},
		},
	})
}

func testAccDataSourceHostConfig_basic(regionID string, hostName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name      = "test-vpc"
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
	name = "test-hostgroup"
}

resource "ecloud_host" "test-host" {
	host_group_id = ecloud_hostgroup.test-hostgroup.id
	name = "%[2]s"
}

data "ecloud_host" "test-host" {
    name = ecloud_host.test-host.name
}
`, regionID, hostName)
}
