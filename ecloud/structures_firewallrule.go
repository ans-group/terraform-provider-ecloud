package ecloud

import "github.com/ukfast/sdk-go/pkg/service/ecloud"

func expandCreateFirewallRuleRequestPorts(rawPorts []interface{}) ([]ecloud.CreateFirewallRulePortRequest, error) {
	var ports []ecloud.CreateFirewallRulePortRequest

	for _, rawPort := range rawPorts {
		port := rawPort.(map[string]interface{})
		protocol := port["protocol"].(string)
		protocolParsed, err := ecloud.ParseFirewallRulePortProtocol(protocol)
		if err != nil {
			return nil, err
		}

		ports = append(ports, ecloud.CreateFirewallRulePortRequest{
			Protocol:    protocolParsed,
			Source:      port["source"].(string),
			Destination: port["destination"].(string),
		})
	}

	return ports, nil
}

func expandUpdateFirewallRuleRequestPorts(rawPorts []interface{}) ([]ecloud.PatchFirewallRulePortRequest, error) {
	var ports []ecloud.PatchFirewallRulePortRequest

	for _, rawPort := range rawPorts {
		port := rawPort.(map[string]interface{})
		protocol := port["protocol"].(string)
		protocolParsed, err := ecloud.ParseFirewallRulePortProtocol(protocol)
		if err != nil {
			return nil, err
		}

		ports = append(ports, ecloud.PatchFirewallRulePortRequest{
			Protocol:    protocolParsed,
			Source:      port["source"].(string),
			Destination: port["destination"].(string),
		})
	}

	return ports, nil
}
