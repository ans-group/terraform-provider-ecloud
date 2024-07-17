package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRouter_basic(t *testing.T) {
	routerName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_router.test-router"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRouterConfig_basic(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", routerName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckRouterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Router ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetRouter(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.RouterNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckRouterDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_router" {
			continue
		}

		_, err := service.GetRouter(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Router with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.RouterNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceRouterConfig_basic(routerName string) string {
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

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%s"
}
`, routerName)
}
