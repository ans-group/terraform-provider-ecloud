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
					resource.TestCheckResourceAttr(resourceName, "vpc_id", UKF_TEST_REFERENCE_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_id", UKF_TEST_REFERENCE_DHCP_AVAILABILITY_ZONE_ID),
				),
			},
		},
	})
}

var testAccCheckDataSourceDHCPConfigTemplate = `
data "ecloud_dhcp" "test-dhcp" {
    dhcp_id = "{{ .UKF_TEST_REFERENCE_DHCP_ID }}"
}`

func testAccCheckDataSourceDHCPConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_REFERENCE_DHCP_ID": UKF_TEST_REFERENCE_DHCP_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceDHCPConfigTemplate, data)
}
