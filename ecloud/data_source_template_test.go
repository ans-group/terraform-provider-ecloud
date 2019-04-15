package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceTemplate(t *testing.T) {
	var template ecloudservice.Template
	config, err := testAccCheckDataSourceTemplateConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_template.test-template"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceTemplateExists(resourceName, &template),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_TEMPLATE_NAME),
					resource.TestCheckResourceAttr(resourceName, "platform", UKF_TEST_TEMPLATE_PLATFORM),
					resource.TestCheckResourceAttr(resourceName, "pod_id", UKF_TEST_SOLUTION_POD_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourceTemplateExists(n string, template *ecloudservice.Template) resource.TestCheckFunc {
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

var testAccCheckDataSourceTemplateConfigTemplate = `
data "ecloud_template" "test-template" {
	name = "{{ .UKF_TEST_TEMPLATE_NAME }}"
	pod_id = "{{ .UKF_TEST_SOLUTION_POD_ID }}"
}`

func testAccCheckDataSourceTemplateConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_TEMPLATE_NAME":   UKF_TEST_TEMPLATE_NAME,
		"UKF_TEST_SOLUTION_POD_ID": UKF_TEST_SOLUTION_POD_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceTemplateConfigTemplate, data)
}
