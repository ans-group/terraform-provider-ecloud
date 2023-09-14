package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHost_basic(t *testing.T) {
	hostName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_host.test-host"
	hostGroupResourceName := "ecloud_hostgroup.test-hostgroup"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHostConfig_basic(ANS_TEST_VPC_REGION_ID, hostName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", hostName),
					resource.TestCheckResourceAttrPair(hostGroupResourceName, "id", resourceName, "host_group_id"),
				),
			},
		},
	})
}

func testAccCheckHostExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No host ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetHost(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.HostNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckHostDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_host" {
			continue
		}

		_, err := service.GetHost(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Host with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.HostNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceHostConfig_basic(regionID string, HostName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_hostspec" "test-hostspec" {
	name = "DUAL-4208--64GB"
}

resource "ecloud_hostgroup" "test-hostgroup" {
	vpc_id = ecloud_vpc.test-vpc.id
	host_spec_id = data.ecloud_hostspec.test-hostspec.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-hostgroup"
	windows_enabled = false
}

resource "ecloud_host" "test-host" {
	host_group_id = ecloud_hostgroup.test-hostgroup.id
	name = "%[2]s"
}
`, regionID, HostName)
}
