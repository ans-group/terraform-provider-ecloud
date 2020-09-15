package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceNetwork(t *testing.T) {
	var network ecloudservice.Network
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
					testAccCheckDataSourceNetworkExists(resourceName, &network),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_NETWORK_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceNetworkExists(n string, network *ecloudservice.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getNetwork, err := service.GetNetwork(rs.Primary.ID)
		if err != nil {
			return err
		}

		*network = getNetwork

		return nil
	}
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
