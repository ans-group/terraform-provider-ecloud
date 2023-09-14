package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRouterThroughput_basic(t *testing.T) {
	throughputName := "1GB"
	config := testAccDataSourceRouterThroughputConfig_basic(throughputName)
	resourceName := "data.ecloud_router_throughput.test-throughput"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", throughputName),
				),
			},
		},
	})
}

func testAccDataSourceRouterThroughputConfig_basic(routerName string) string {
	return fmt.Sprintf(`
data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_router_throughput" "test-throughput" {
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%s"
}
`, routerName)
}
