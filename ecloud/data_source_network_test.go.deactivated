package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNetwork(t *testing.T) {
	config, err := testAccCheckDataSourceNetworkConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_network.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_REFERENCE_NETWORK_NAME),
					resource.TestCheckResourceAttr(resourceName, "router_id", UKF_TEST_REFERENCE_ROUTER_ID),
				),
			},
		},
	})
}

var testAccCheckDataSourceNetworkConfigTemplate = `
data "ecloud_network" "test-network" {
    network_id = "{{ .UKF_TEST_REFERENCE_NETWORK_ID }}"
}`

func testAccCheckDataSourceNetworkConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_REFERENCE_NETWORK_ID": UKF_TEST_REFERENCE_NETWORK_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceNetworkConfigTemplate, data)
}
