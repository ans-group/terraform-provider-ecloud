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

func TestAccSolutionTemplate_basic(t *testing.T) {
	var template ecloudservice.Template
	vmName := acctest.RandomWithPrefix("tftest")
	templateName := acctest.RandomWithPrefix("tftest")
	config, err := testAccCheckSolutionTemplateConfig_basic(vmName, templateName)
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "ecloud_solution_template.test-template"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolutionTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolutionTemplateExists(resourceName, &template),
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
				),
			},
		},
	})
}

func testAccCheckSolutionTemplateExists(n string, template *ecloudservice.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		templateName := rs.Primary.ID

		solutionID, err := strconv.Atoi(rs.Primary.Attributes["solution_id"])
		if err != nil {
			return err
		}

		getTemplate, err := service.GetSolutionTemplate(solutionID, templateName)
		if err != nil {
			if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
				return nil
			}
			return err
		}

		*template = getTemplate

		return nil
	}
}

func testAccCheckSolutionTemplateDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_solution_template" {
			continue
		}

		templateName := rs.Primary.ID

		solutionID, err := strconv.Atoi(rs.Primary.Attributes["solution_id"])
		if err != nil {
			return err
		}

		_, err = service.GetSolutionTemplate(solutionID, templateName)
		if err == nil {
			return fmt.Errorf("Template with name [%s] for solution with ID [%d] still exists", templateName, solutionID)
		}

		if _, ok := err.(*ecloudservice.TemplateNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckSolutionTemplateConfigTemplate_basic = `
resource "ecloud_virtualmachine" "test-vm" {
  cpu = 2
  ram = 2
  os_disk = 20
  template = "CentOS 7 64-bit"
  name = "{{ .VMName }}"
  environment = "{{ .UKF_TEST_SOLUTION_ENVIRONMENT }}"
  solution_id = {{ .UKF_TEST_SOLUTION_ID }}
}

resource "ecloud_solution_template" "test-template" {
	virtualmachine_id = "${ecloud_virtualmachine.test-vm.id}"
	name = "{{ .TemplateName }}"
	solution_id = {{ .UKF_TEST_SOLUTION_ID }}
}`

func testAccCheckSolutionTemplateConfig_basic(vmName string, templateName string) (string, error) {
	data := map[string]interface{}{
		"VMName":                        vmName,
		"UKF_TEST_SOLUTION_ENVIRONMENT": UKF_TEST_SOLUTION_ENVIRONMENT,
		"UKF_TEST_SOLUTION_ID":          UKF_TEST_SOLUTION_ID,
		"TemplateName":                  templateName,
	}

	return testAccTemplateConfig(testAccCheckSolutionTemplateConfigTemplate_basic, data)
}
