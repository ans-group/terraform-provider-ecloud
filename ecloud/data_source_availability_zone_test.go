package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceAvailabilityZone(t *testing.T) {
	var az ecloudservice.AvailabilityZone
	config, err := testAccCheckDataSourceAvailabilityZoneConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_availabilityzone.test-az"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceAvailabilityZoneExists(resourceName, &az),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_AVAILABILITYZONE_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceAvailabilityZoneExists(n string, az *ecloudservice.AvailabilityZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No availability zone ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getAvailabilityZone, err := service.GetAvailabilityZone(rs.Primary.ID)
		if err != nil {
			return err
		}

		*az = getAvailabilityZone

		return nil
	}
}

var testAccCheckDataSourceAvailabilityZoneConfigTemplate = `
data "ecloud_availabilityzone" "test-az" {
    name = "{{ .UKF_TEST_AVAILABILITYZONE_NAME }}"
}`

func testAccCheckDataSourceAvailabilityZoneConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_AVAILABILITYZONE_NAME": UKF_TEST_AVAILABILITYZONE_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceAvailabilityZoneConfigTemplate, data)
}
