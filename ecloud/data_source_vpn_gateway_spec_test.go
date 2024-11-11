package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNGatewaySpecification_basic(t *testing.T) {
	config := testAccDataSourceVPNGatewaySpecificationConfig_basic("Small")
	resourceName := "data.ecloud_vpn_gateway_specification.test-vpngatewayspec"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Small"),
				),
			},
		},
	})
}

func testAccDataSourceVPNGatewaySpecificationConfig_basic(specName string) string {
	return fmt.Sprintf(`
data "ecloud_vpn_gateway_specification" "test-vpngatewayspec" {
	name = "%[1]s"
}
`, specName)
}
