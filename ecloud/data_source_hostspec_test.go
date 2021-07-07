package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceHostSpec_basic(t *testing.T) {
	hostSpecName := "DUAL-E5-2620--32GB"
	config := testAccDataSourceHostSpecConfig_basic(hostSpecName)
	resourceName := "data.ecloud_hostspec.test-hostspec"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
