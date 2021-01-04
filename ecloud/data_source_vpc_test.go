package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVPC_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPCConfig_basic(UKF_TEST_VPC_REGION_ID, vpcName)
	resourceName := "data.ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region_id", UKF_TEST_VPC_REGION_ID),
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccDataSourceVPCConfig_basic(regionID string, vpcName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%s"
    name = "%s"
}

data "ecloud_vpc" "test-vpc" {
    name = ecloud_vpc.test-vpc.name
}
`, regionID, vpcName)
}
