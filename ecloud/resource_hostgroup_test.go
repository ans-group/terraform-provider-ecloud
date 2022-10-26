package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHostGroup_basic(t *testing.T) {
	hostGroupName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_hostgroup.test-hostgroup"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHostGroupConfig_basic(ANS_TEST_VPC_REGION_ID, hostGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", hostGroupName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckHostGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No host group ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetHostGroup(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.HostGroupNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckHostGroupDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_hostgroup" {
			continue
		}

		_, err := service.GetHostGroup(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Host group with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.HostGroupNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceHostGroupConfig_basic(regionID string, hostGroupName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_hostspec" "test-hostspec" {
	name = "DUAL-E5-2620--32GB"
}

resource "ecloud_hostgroup" "test-hostgroup" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	host_spec_id = data.ecloud_hostspec.test-hostspec.id
	name = "%[2]s"
	windows_enabled = false
}
`, regionID, hostGroupName)
}
