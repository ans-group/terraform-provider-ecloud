package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTag_basic(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceTagConfig_basic(tagName)
	resourceName := "data.ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
				),
			},
		},
	})
}

func TestAccDataSourceTag_withScope(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	tagScope := "vpc"
	config := testAccDataSourceTagConfig_withScope(tagName, tagScope)
	resourceName := "data.ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
					resource.TestCheckResourceAttr(resourceName, "scope", tagScope),
				),
			},
		},
	})
}

func TestAccDataSourceTag_byID(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceTagConfig_byID(tagName)
	resourceName := "data.ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
				),
			},
		},
	})
}

func testAccDataSourceTagConfig_basic(tagName string) string {
	return fmt.Sprintf(`
resource "ecloud_tag" "test-tag" {
	name = "%s"
}

data "ecloud_tag" "test-tag" {
    name = ecloud_tag.test-tag.name
}
`, tagName)
}

func testAccDataSourceTagConfig_withScope(tagName, tagScope string) string {
	return fmt.Sprintf(`
resource "ecloud_tag" "test-tag" {
	name = "%s"
	scope = "%s"
}

data "ecloud_tag" "test-tag" {
    name = ecloud_tag.test-tag.name
    scope = ecloud_tag.test-tag.scope
}
`, tagName, tagScope)
}

func testAccDataSourceTagConfig_byID(tagName string) string {
	return fmt.Sprintf(`
resource "ecloud_tag" "test-tag" {
	name = "%s"
}

data "ecloud_tag" "test-tag" {
    tag_id = ecloud_tag.test-tag.id
}
`, tagName)
}
