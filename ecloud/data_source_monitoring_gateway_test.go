package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMonitoringGateway_basic(t *testing.T) {
	gatewayName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceMonitoringGatewayConfig_basic(gatewayName)
	resourceName := "data.ecloud_monitoring_gateway.test-gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", gatewayName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttrSet(resourceName, "specification_id"),
				),
			},
		},
	})
}

func testAccDataSourceMonitoringGatewayConfig_basic(gatewayName string) string {
	return fmt.Sprintf(`
data "ecloud_monitoring_gateway" "test-gateway" {
	name = "%s
}
`, gatewayName)
}
