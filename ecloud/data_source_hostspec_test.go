package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHostSpec_basic(t *testing.T) {
	hostSpecName := "DUAL-4208--64GB"
	config := testAccDataSourceHostSpecConfig_basic(hostSpecName)
	resourceName := "data.ecloud_hostspec.test-hostspec"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", hostSpecName),
				),
			},
		},
	})
}

func testAccDataSourceHostSpecConfig_basic(hostSpecName string) string {
	return fmt.Sprintf(`

data "ecloud_hostspec" "test-hostspec" {
    name = "%s"
}
`, hostSpecName)
}
