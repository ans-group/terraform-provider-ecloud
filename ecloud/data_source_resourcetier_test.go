package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceResourceTier_basic(t *testing.T) {
	config := testAccDataSourceResourceTierConfig_basic(ANS_TEST_VPC_REGION_ID)
	resourceName := "data.ecloud_resourcetier.test-rt"
	tierName := "Standard CPU"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tierName),
				),
			},
		},
	})
}

func testAccDataSourceResourceTierConfig_basic(tierName string) string {
	return fmt.Sprintf(`
data "ecloud_resourcetier" "test-rt" {
    name = "%s"
	availability_zone_id = az-4c31a488
}
`, tierName)
}
