package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAvailabilityZone(t *testing.T) {

	azID := "az-4fcc2a10"
	azName := "London Central"
	azCode := "lon1"
	azDCSiteID := 4

	config, err := testAccCheckDataSourceAvailabilityZoneConfig(azID)
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_availability_zone.test-az"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", azName),
					resource.TestCheckResourceAttr(resourceName, "code", azCode),
					resource.TestCheckResourceAttr(resourceName, "datacentre_site_id", fmt.Sprintf("%d", azDCSiteID)),
				),
			},
		},
	})
}

var testAccCheckDataSourceAvailabilityZoneConfigTemplate = `
data "ecloud_availability_zone" "test-az" {
    availability_zone_id = "{{ .UKF_TEST_AVAILABILITYZONE_ID }}"
}`

func testAccCheckDataSourceAvailabilityZoneConfig(azID string) (string, error) {
	return testAccTemplateConfig(testAccCheckDataSourceAvailabilityZoneConfigTemplate, map[string]interface{}{
		"UKF_TEST_AVAILABILITYZONE_ID": azID,
	})
}
