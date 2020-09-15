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

	t.Log(config)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_NETWORK_NAME),
				),
			},
		},
	})
}

var testAccCheckDataSourceNetworkConfigTemplate = `
data "ecloud_network" "test-network" {
    name = "{{ .UKF_TEST_NETWORK_NAME }}"
}`

func testAccCheckDataSourceNetworkConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_NETWORK_NAME": UKF_TEST_NETWORK_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceNetworkConfigTemplate, data)
}
