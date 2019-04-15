package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceSolutionTemplate(t *testing.T) {
	var template ecloudservice.Template
	config, err := testAccCheckDataSourceSolutionTemplateConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_solution_template.test-template"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSolutionTemplateExists(resourceName, &template),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_SOLUTION_TEMPLATE_NAME),
					resource.TestCheckResourceAttr(resourceName, "platform", UKF_TEST_SOLUTION_TEMPLATE_PLATFORM),
					resource.TestCheckResourceAttr(resourceName, "solution_id", UKF_TEST_SOLUTION_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourceSolutionTemplateExists(n string, template *ecloudservice.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		return nil
	}
}

var testAccCheckDataSourceSolutionTemplateConfigTemplate = `
data "ecloud_solution_template" "test-template" {
	name = "{{ .UKF_TEST_SOLUTION_TEMPLATE_NAME }}"
	solution_id = "{{ .UKF_TEST_SOLUTION_ID }}"
}`

func testAccCheckDataSourceSolutionTemplateConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_SOLUTION_TEMPLATE_NAME": UKF_TEST_SOLUTION_TEMPLATE_NAME,
		"UKF_TEST_SOLUTION_ID":            UKF_TEST_SOLUTION_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceSolutionTemplateConfigTemplate, data)
}
