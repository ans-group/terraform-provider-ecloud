package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIPAddress_basic(t *testing.T) {
	params := map[string]string{
		"datasource_name": "test-ipaddress",
		"name":            acctest.RandomWithPrefix("tftest"),
		"ip_address":      "10.0.0.10",
	}
	config := testAccDataSourceIPAddressConfig_basic(params)
	resourceName := "data.ecloud_ipaddress." + params["datasource_name"]

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["name"]),
					resource.TestCheckResourceAttr(resourceName, "ip_address", params["ip_address"]),
				),
			},
		},
	})
}

func testAccDataSourceIPAddressConfig_basic(params map[string]string) string {
	str, _ := testAccTemplateConfig(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_ipaddress" "test-ipaddress" {
	network_id = ecloud_network.test-network.id
	name = "{{ .name }}"
	ip_address = "{{ .ip_address }}"
}

data "ecloud_ipaddress" "{{ .datasource_name }}" {
	ip_address_id = ecloud_ipaddress.test-ipaddress.id
}
`, params)
	return str
}
