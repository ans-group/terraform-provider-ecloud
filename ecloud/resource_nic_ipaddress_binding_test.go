package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccIPAddressNICBinding_basic(t *testing.T) {
	ipAddressName := acctest.RandomWithPrefix("tftest")
	ipAddressIPAddressNICBinding := "10.0.0.5"
	resourceName := "ecloud_ipaddress.test-ipaddress"
	networkResourceName := "ecloud_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAddressNICBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIPAddressNICBindingConfig_basic(UKF_TEST_VPC_REGION_ID, ipAddressName, ipAddressIPAddressNICBinding),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAddressNICBindingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", ipAddressName),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddressIPAddressNICBinding),
					resource.TestCheckResourceAttrPair(networkResourceName, "id", resourceName, "network_id"),
				),
			},
		},
	})
}

func testAccCheckIPAddressNICBindingExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetIPAddressNICBinding(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.IPAddressNICBindingNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckIPAddressNICBindingDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_ipaddress" {
			continue
		}

		_, err := service.GetIPAddressNICBinding(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("IP address with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.IPAddressNICBindingNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceIPAddressNICBindingConfig_basic(regionID string, ipAddressName string, ipAddressIPAddressNICBinding string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_ipaddress" "test-host" {
	network_id = ecloud_network.test-network.id
	name = "%[2]s"
	ip_address = "%[3]s"
}
`, regionID, ipAddressName, ipAddressIPAddressNICBinding)
}
