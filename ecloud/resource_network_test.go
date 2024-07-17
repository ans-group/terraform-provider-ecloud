package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetwork_basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_network.test-network"
	routerResourceName := "ecloud_router.test-router"
	subnet := "10.0.0.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkConfig_basic(networkName, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", networkName),
					resource.TestCheckResourceAttrPair(routerResourceName, "id", resourceName, "router_id"),
					resource.TestCheckResourceAttr(resourceName, "subnet", subnet),
				),
			},
		},
	})
}

func testAccCheckNetworkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetNetwork(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NetworkNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckNetworkDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_network" {
			continue
		}

		_, err := service.GetNetwork(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.NetworkNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceNetworkConfig_basic(networkName string, subnet string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
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
	name = "%[1]s"
	subnet = "%[2]s"
}
`, networkName, subnet)
}
