package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceRouterThroughput_basic(t *testing.T) {
	throughputName := "1GB"
	config := testAccDataSourceRouterThroughputConfig_basic(throughputName)
	resourceName := "data.ecloud_router_throughput.test-throughput"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
data "ecloud_router_throughput" "test-throughput" {
	# Hard code AZ ID whilst we currently do not expose via API
	availability_zone_id = "az-4fcc2a10"
    name = "%s"
}
`, routerName)
}
