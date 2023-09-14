package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIOPS_basic(t *testing.T) {
	level := 300
	config := testAccDataSourceIOPSConfig_basic(level)
	resourceName := "data.ecloud_iops.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "level", strconv.Itoa(level)),
				),
			},
		},
	})
}

func testAccDataSourceIOPSConfig_basic(level int) string {
	return fmt.Sprintf(`
data "ecloud_iops" "test" {
    level = "%d"
}
`, level)
}
