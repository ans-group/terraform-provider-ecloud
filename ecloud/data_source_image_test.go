package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceImage_basic(t *testing.T) {
	imageName := "CentOS 7"
	config := testAccDataSourceImageConfig_basic(imageName)
	resourceName := "data.ecloud_image.test-image"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", imageName),
				),
			},
		},
	})
}

func testAccDataSourceImageConfig_basic(imageName string) string {
	return fmt.Sprintf(`
data "ecloud_image" "test-image" {
    name = "%s"
}
`, imageName)
}
