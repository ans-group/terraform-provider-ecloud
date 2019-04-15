package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceSolution(t *testing.T) {
	var solution ecloudservice.Solution
	config, err := testAccCheckDataSourceSolutionConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_solution.test-solution"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSolutionExists(resourceName, &solution),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_SOLUTION_NAME),
					resource.TestCheckResourceAttr(resourceName, "environment", UKF_TEST_SOLUTION_ENVIRONMENT),
				),
			},
		},
	})
}

func testAccCheckDataSourceSolutionExists(n string, solution *ecloudservice.Solution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No solution ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		solutionID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		getSolution, err := service.GetSolution(solutionID)
		if err != nil {
			return err
		}

		*solution = getSolution

		return nil
	}
}

var testAccCheckDataSourceSolutionConfigTemplate = `
data "ecloud_solution" "test-solution" {
	name = "{{ .UKF_TEST_SOLUTION_NAME }}"
}`

func testAccCheckDataSourceSolutionConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_SOLUTION_NAME": UKF_TEST_SOLUTION_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceSolutionConfigTemplate, data)
}
