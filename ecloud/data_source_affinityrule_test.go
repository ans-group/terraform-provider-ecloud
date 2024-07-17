package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAffinityRule_basic(t *testing.T) {
	arName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceAffinityRuleConfig_basic(arName)
	resourceName := "data.ecloud_affinityrule.test-ar"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", arName),
				),
			},
		},
	})
}

func testAccDataSourceAffinityRuleConfig_basic(arName string) string {
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

resource "ecloud_affinityrule" "test-ar" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%[1]s"
	type = "anti-affinity"
}

data "ecloud_affinityrule" "test-ar" {
    name = ecloud_affinityrule.test-ar.name
}
`, arName)
}
