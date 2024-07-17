package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPC_basic(t *testing.T) {
	vpcName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceVPCConfig_basic(vpcName)
	resourceName := "data.ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", vpcName),
				),
			},
		},
	})
}

func testAccDataSourceVPCConfig_basic(vpcName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
    name = "%s"
}

data "ecloud_vpc" "test-vpc" {
    name = ecloud_vpc.test-vpc.name
}
`, vpcName)
}
