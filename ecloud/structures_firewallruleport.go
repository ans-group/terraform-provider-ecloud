package ecloud

import "github.com/ukfast/sdk-go/pkg/service/ecloud"

func flattenFirewallRulePorts(ports []ecloud.FirewallRulePort) interface{} {
	var flattenedPorts []map[string]interface{}

	for _, port := range ports {
		flattenedPorts = append(flattenedPorts, map[string]interface{}{
			"name":             port.Name,
			"firewall_rule_id": port.FirewallRuleID,
			"protocol":         port.Protocol,
			"source":           port.Source,
			"destination":      port.Destination,
		})
	}

	return flattenedPorts
}
