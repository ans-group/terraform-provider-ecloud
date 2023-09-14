package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLoadBalancerSpec_basic(t *testing.T) {
	lbSpecName := "Medium"
	config := testAccDataSourceLoadBalancerSpecConfig_basic(lbSpecName)
	resourceName := "data.ecloud_loadbalancer_spec.test-lbspec"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", lbSpecName),
				),
			},
		},
	})
}

func testAccDataSourceLoadBalancerSpecConfig_basic(lbSpecName string) string {
	return fmt.Sprintf(`

data "ecloud_loadbalancer_spec" "test-lbspec" {
    name = "%s"
}
`, lbSpecName)
}
