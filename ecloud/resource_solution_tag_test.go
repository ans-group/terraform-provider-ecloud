package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccSolutionTag_basic(t *testing.T) {
	var tag ecloudservice.Tag
	tagKey := acctest.RandomWithPrefix("tftest")
	config, err := testAccCheckSolutionTagConfig_basic(tagKey, "bar")
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "ecloud_solution_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolutionTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolutionTagExists(resourceName, &tag),
					resource.TestCheckResourceAttr(resourceName, "key", tagKey),
					resource.TestCheckResourceAttr(resourceName, "value", "bar"),
				),
			},
		},
	})
}

func testAccCheckSolutionTagExists(n string, tag *ecloudservice.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No solution tag ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		tagKey := rs.Primary.ID

		solutionID, err := strconv.Atoi(rs.Primary.Attributes["solution_id"])
		if err != nil {
			return err
		}

		getTag, err := service.GetSolutionTag(solutionID, tagKey)
		if err != nil {
			if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
				return nil
			}
			return err
		}

		*tag = getTag

		return nil
	}
}

func testAccCheckSolutionTagDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_solution_tag" {
			continue
		}

		tagKey := rs.Primary.ID

		solutionID, err := strconv.Atoi(rs.Primary.Attributes["solution_id"])
		if err != nil {
			return err
		}

		_, err = service.GetSolutionTag(solutionID, tagKey)
		if err == nil {
			return fmt.Errorf("Tag with key [%s] for solution with ID [%d] still exists", tagKey, solutionID)
		}

		if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckSolutionTagConfigTemplate_basic = `
resource "ecloud_solution_tag" "test-tag" {
	solution_id = "{{ .UKF_TEST_SOLUTION_ID }}"
    key = "{{ .TagKey }}"
    value = "{{ .TagValue }}"
}`

func testAccCheckSolutionTagConfig_basic(tagKey string, tagValue string) (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_SOLUTION_ID": UKF_TEST_SOLUTION_ID,
		"TagKey":               tagKey,
		"TagValue":             tagValue,
	}

	return testAccTemplateConfig(testAccCheckSolutionTagConfigTemplate_basic, data)
}
