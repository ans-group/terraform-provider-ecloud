package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourcePodAppliance(t *testing.T) {
	var appliance ecloudservice.Appliance
	config, err := testAccCheckDataSourcePodApplianceConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_pod_appliance.test-pod-appliance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourcePodApplianceExists(resourceName, &appliance),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_APPLIANCE_NAME),
					resource.TestCheckResourceAttr(resourceName, "pod_id", UKF_TEST_APPLIANCE_POD_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourcePodApplianceExists(n string, appliance *ecloudservice.Appliance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No appliance ID is set")
		}

		return nil
	}
}

var testAccCheckDataSourcePodApplianceConfigTemplate = `
data "ecloud_pod_appliance" "test-pod-appliance" {
	name = "{{ .UKF_TEST_APPLIANCE_NAME }}"
	pod_id = "{{ .UKF_TEST_APPLIANCE_POD_ID }}"
}`

func testAccCheckDataSourcePodApplianceConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_APPLIANCE_NAME":   UKF_TEST_APPLIANCE_NAME,
		"UKF_TEST_APPLIANCE_POD_ID": UKF_TEST_SOLUTION_POD_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourcePodApplianceConfigTemplate, data)
}
