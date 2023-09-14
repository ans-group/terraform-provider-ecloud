package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkRule_basic(t *testing.T) {
	params := map[string]string{
		"vpc_region_id":    ANS_TEST_VPC_REGION_ID,
		"rule_name":        acctest.RandomWithPrefix("tftest"),
		"rule_sequence":    "0",
		"rule_direction":   "IN",
		"rule_action":      "ALLOW",
		"rule_source":      "10.0.0.5/32",
		"rule_destination": "ANY",
		"rule_enabled":     "true",
	}
	resourceName := "ecloud_networkrule.test-nr"
	policyResourceName := "ecloud_networkpolicy.test-np"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetworkRuleConfig_basic(params),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkRuleExists(resourceName),
					resource.TestCheckResourceAttrPair(policyResourceName, "id", resourceName, "network_policy_id"),
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

func testAccCheckNetworkRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network rule ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetNetworkRule(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NetworkRuleNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckNetworkRuleDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_networkrule" {
			continue
		}

		_, err := service.GetNetworkRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network rule with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.NetworkRuleNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceNetworkRuleConfig_basic(params map[string]string) string {
	str, _ := testAccTemplateConfig(`
	resource "ecloud_vpc" "test-vpc" {
		region_id = "{{ .vpc_region_id }}"
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
	`, params)

	return str
}
