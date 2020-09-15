package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceDHCP(t *testing.T) {
	var dhcp ecloudservice.DHCP
	config, err := testAccCheckDataSourceDHCPConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_dhcp.test-dhcp"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDHCPExists(resourceName, &dhcp),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_id", UKF_TEST_DHCP_AVAILABILITYZONE_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourceDHCPExists(n string, dhcp *ecloudservice.DHCP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getDHCP, err := service.GetDHCP(rs.Primary.ID)
		if err != nil {
			return err
		}

		*dhcp = getDHCP

		return nil
	}
}

var testAccCheckDataSourceDHCPConfigTemplate = `
data "ecloud_dhcp" "test-dhcp" {
    availability_zone_id = "{{ .UKF_TEST_DHCP_AVAILABILITYZONE_ID }}"
}`

func testAccCheckDataSourceDHCPConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_DHCP_AVAILABILITYZONE_ID": UKF_TEST_DHCP_AVAILABILITYZONE_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceDHCPConfigTemplate, data)
}
