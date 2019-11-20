package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceAppliance(t *testing.T) {
	var appliance ecloudservice.Appliance
	config, err := testAccCheckDataSourceApplianceConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_appliance.test-appliance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceApplianceExists(resourceName, &appliance),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_APPLIANCE_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceApplianceExists(n string, appliance *ecloudservice.Appliance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No appliance ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		var applianceID = rs.Primary.ID

		getAppliance, err := service.GetAppliance(applianceID)
		if err != nil {
			return err
		}

		*appliance = getAppliance

		return nil
	}
}

var testAccCheckDataSourceApplianceConfigTemplate = `
data "ecloud_appliance" "test-appliance" {
    name = "{{ .UKF_TEST_APPLIANCE_NAME }}"
}`

func testAccCheckDataSourceApplianceConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_APPLIANCE_NAME": UKF_TEST_APPLIANCE_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceApplianceConfigTemplate, data)
}
