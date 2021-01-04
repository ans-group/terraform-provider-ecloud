package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccRouter_basic(t *testing.T) {
	routerName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_router.test-router"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRouterConfig_basic(UKF_TEST_VPC_REGION_ID, routerName),
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

func testAccResourceRouterConfig_basic(regionID string, routerName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "%s"
}
`, regionID, routerName)
}
