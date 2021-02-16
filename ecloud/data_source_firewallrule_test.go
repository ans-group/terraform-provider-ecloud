package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceFirewallRule_basic(t *testing.T) {
	params := map[string]string{
		"vpc_region_id":    UKF_TEST_VPC_REGION_ID,
		"rule_name":        acctest.RandomWithPrefix("tftest"),
		"rule_sequence":    "0",
		"rule_direction":   "IN",
		"rule_action":      "ALLOW",
		"rule_source":      "192.168.1.1/32",
		"rule_destination": "ANY",
		"rule_enabled":     "true",
	}
	config := testAccDataSourceFirewallRuleConfig_basic(params)
	resourceName := "data.ecloud_firewallrule.test-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func testAccDataSourceFirewallRuleConfig_basic(params map[string]string) string {
	str, _ := testAccTemplateConfig(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "{{ .vpc_region_id }}"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_firewallpolicy" "test-fwp" {
	router_id = ecloud_router.test-router.id
	name = "test-fwp"
	sequence = 0
}

resource "ecloud_firewallrule" "test-rule" {
	firewall_policy_id = ecloud_firewallpolicy.test-fwp.id
	name = "{{ .rule_name }}"
	sequence = {{ .rule_sequence }}
	direction = "{{ .rule_direction }}"
	source = "{{ .rule_source }}"
	destination = "{{ .rule_destination }}"
	action = "{{ .rule_action }}"
	enabled = {{ .rule_enabled }}
}

data "ecloud_firewallrule" "test-rule" {
    name = ecloud_firewallrule.test-rule.name
}
`, params)
	return str
}
