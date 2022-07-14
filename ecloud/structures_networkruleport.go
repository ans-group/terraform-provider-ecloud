package ecloud

import "github.com/ans-group/sdk-go/pkg/service/ecloud"

func flattenNetworkRulePorts(ports []ecloud.NetworkRulePort) interface{} {
	var flattenedPorts []map[string]interface{}

	for _, port := range ports {
		flattenedPorts = append(flattenedPorts, map[string]interface{}{
			"name":            port.Name,
			"network_rule_id": port.NetworkRuleID,
			"protocol":        port.Protocol,
			"source":          port.Source,
			"destination":     port.Destination,
		})
	}

	return flattenedPorts
}
