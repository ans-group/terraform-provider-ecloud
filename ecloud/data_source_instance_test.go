package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceInstance(t *testing.T) {
	var instance ecloudservice.Instance
	config, err := testAccCheckDataSourceInstanceConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_instance.test-instance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_INSTANCE_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceInstanceExists(n string, instance *ecloudservice.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getInstance, err := service.GetInstance(rs.Primary.ID)
		if err != nil {
			return err
		}

		*instance = getInstance

		return nil
	}
}

var testAccCheckDataSourceInstanceConfigTemplate = `
data "ecloud_instance" "test-instance" {
    name = "{{ .UKF_TEST_INSTANCE_NAME }}"
}`

func testAccCheckDataSourceInstanceConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_INSTANCE_NAME": UKF_TEST_INSTANCE_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceInstanceConfigTemplate, data)
}
