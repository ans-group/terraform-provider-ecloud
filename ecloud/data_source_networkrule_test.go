package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkRule_basic(t *testing.T) {
	params := map[string]string{
		"rule_name":        acctest.RandomWithPrefix("tftest"),
		"rule_sequence":    "0",
		"rule_direction":   "IN",
		"rule_action":      "ALLOW",
		"rule_source":      "10.0.0.5/32",
		"rule_destination": "ANY",
		"rule_enabled":     "true",
	}
	config := testAccDataSourceNetworkRuleConfig_basic(params)
	resourceName := "data.ecloud_networkrule.test-nr"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["rule_name"]),
					resource.TestCheckResourceAttr(resourceName, "sequence", params["rule_sequence"]),
					resource.TestCheckResourceAttr(resourceName, "direction", params["rule_direction"]),
					resource.TestCheckResourceAttr(resourceName, "action", params["rule_action"]),
					resource.TestCheckResourceAttr(resourceName, "source", params["rule_source"]),
					resource.TestCheckResourceAttr(resourceName, "destination", params["rule_destination"]),
					resource.TestCheckResourceAttr(resourceName, "enabled", params["rule_enabled"]),
				),
			},
		},
	})
}

func testAccDataSourceNetworkRuleConfig_basic(params map[string]string) string {
	str, _ := testAccTemplateConfig(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "test-vpc"
	advanced_networking = true
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_networkpolicy" "test-np" {
	network_id = ecloud_network.test-network.id
	name = "test-policy"
	catchall_rule_action = "REJECT"
}

resource "ecloud_networkrule" "test-nr" {
	network_policy_id = ecloud_networkpolicy.test-np.id
	name = "{{ .rule_name }}"
	sequence = {{ .rule_sequence }}
	direction = "{{ .rule_direction }}"
	source = "{{ .rule_source }}"
	destination = "{{ .rule_destination }}"
	action = "{{ .rule_action }}"
	enabled = {{ .rule_enabled }}
}

data "ecloud_networkrule" "test-nr" {
    name = ecloud_networkrule.test-nr.name
}
`, params)
	return str
}
