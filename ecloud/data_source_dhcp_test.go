package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDHCP(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "vpc_id", UKF_TEST_VPC_ID),
				),
			},
		},
	})
}

var testAccCheckDataSourceDHCPConfigTemplate = `
data "ecloud_dhcp" "test-dhcp" {
    availability_zone_id = "{{ .UKF_TEST_DHCP_AVAILABILITY_ZONE_ID }}"
}`

func testAccCheckDataSourceDHCPConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_DHCP_AVAILABILITY_ZONE_ID": UKF_TEST_DHCP_AVAILABILITY_ZONE_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceDHCPConfigTemplate, data)
}
