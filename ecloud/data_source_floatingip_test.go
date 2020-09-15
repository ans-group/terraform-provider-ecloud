package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceFloatingIP(t *testing.T) {
	var az ecloudservice.FloatingIP
	config, err := testAccCheckDataSourceFloatingIPConfig()
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
					testAccCheckDataSourceFloatingIPExists(resourceName, &az),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_AVAILABILITYZONE_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceFloatingIPExists(n string, az *ecloudservice.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No availability zone ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getFloatingIP, err := service.GetFloatingIP(rs.Primary.ID)
		if err != nil {
			return err
		}

		*az = getFloatingIP

		return nil
	}
}

var testAccCheckDataSourceFloatingIPConfigTemplate = `
data "ecloud_availabilityzone" "test-az" {
    name = "{{ .UKF_TEST_AVAILABILITYZONE_NAME }}"
}`

func testAccCheckDataSourceFloatingIPConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_AVAILABILITYZONE_NAME": UKF_TEST_AVAILABILITYZONE_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceFloatingIPConfigTemplate, data)
}
