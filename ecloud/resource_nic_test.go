package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNIC_basic(t *testing.T) {
	nicName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_nic.test-nic"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNICConfig_basic(nicName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNICExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nicName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func testAccCheckNICExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No nic ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetNIC(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NICNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckNICDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_nic" {
			continue
		}

		_, err := service.GetNIC(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("NIC with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.NICNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceNICConfig_basic(nicName string) string {
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

resource "ecloud_nic" "test-nic" {
  name = "%[1]s"
  instance_id = ecloud_instance.instance-1.id
  network_id = ecloud_network.network-1.id
}
`, nicName)
}
