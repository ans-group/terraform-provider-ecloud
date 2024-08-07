package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewallRule_basic(t *testing.T) {
	params := map[string]string{
		"rule_name":        acctest.RandomWithPrefix("tftest"),
		"rule_sequence":    "0",
		"rule_direction":   "IN",
		"rule_action":      "ALLOW",
		"rule_source":      "192.168.1.1/32",
		"rule_destination": "ANY",
		"rule_enabled":     "true",
	}
	resourceName := "ecloud_firewallrule.test-fwr"
	policyResourceName := "ecloud_firewallpolicy.test-fwp"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFirewallRuleConfig_basic(params),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttrPair(policyResourceName, "id", resourceName, "firewall_policy_id"),
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

func testAccCheckFirewallRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall rule ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetFirewallRule(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.FirewallRuleNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckFirewallRuleDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_firewallrule" {
			continue
		}

		_, err := service.GetFirewallRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Firewall rule with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.FirewallRuleNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceFirewallRuleConfig_basic(params map[string]string) string {
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

resource "ecloud_firewallpolicy" "test-fwp" {
	router_id = ecloud_router.test-router.id
	sequence = 0
}

resource "ecloud_firewallrule" "test-fwr" {
	firewall_policy_id = ecloud_firewallpolicy.test-fwp.id
	name = "{{ .rule_name }}"
	sequence = {{ .rule_sequence }}
	direction = "{{ .rule_direction }}"
	source = "{{ .rule_source }}"
	destination = "{{ .rule_destination }}"
	action = "{{ .rule_action }}"
	enabled = {{ .rule_enabled }}
}
`, params)

	return str
}
