package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTag_basic(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
				),
			},
		},
	})
}

func TestAccTag_withScope(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	tagScope := "vpc"
	resourceName := "ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTagConfig_withScope(tagName, tagScope),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
					resource.TestCheckResourceAttr(resourceName, "scope", tagScope),
				),
			},
		},
	})
}

func TestAccTag_update(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	tagNameUpdated := acctest.RandomWithPrefix("tftest-updated")
	resourceName := "ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
				),
			},
			{
				Config: testAccResourceTagConfig_basic(tagNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagNameUpdated),
				),
			},
		},
	})
}

func TestAccTag_updateScope(t *testing.T) {
	tagName := acctest.RandomWithPrefix("tftest")
	tagScope := "vpc"
	tagScopeUpdated := "instance"
	resourceName := "ecloud_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTagConfig_withScope(tagName, tagScope),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
					resource.TestCheckResourceAttr(resourceName, "scope", tagScope),
				),
			},
			{
				Config: testAccResourceTagConfig_withScope(tagName, tagScopeUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", tagName),
					resource.TestCheckResourceAttr(resourceName, "scope", tagScopeUpdated),
				),
			},
		},
	})
}

func testAccCheckTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ecloud: not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ecloud: no Tag ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetTag(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
				return fmt.Errorf("ecloud: tag with ID [%s] not found", rs.Primary.ID)
			}
			return err
		}

		return nil
	}
}

func testAccCheckTagDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_tag" {
			continue
		}

		_, err := service.GetTag(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ecloud: tag with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
			continue
		}

		return err
	}

	return nil
}

func testAccResourceTagConfig_basic(tagName string) string {
	return fmt.Sprintf(`
resource "ecloud_tag" "test-tag" {
	name = "%s"
}
`, tagName)
}

func testAccResourceTagConfig_withScope(tagName, tagScope string) string {
	return fmt.Sprintf(`
resource "ecloud_tag" "test-tag" {
	name = "%s"
	scope = "%s"
}
`, tagName, tagScope)
}
