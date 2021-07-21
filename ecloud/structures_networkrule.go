package ecloud

import "github.com/ukfast/sdk-go/pkg/service/ecloud"

func expandCreateNetworkRuleRequestPorts(rawPorts []interface{}) ([]ecloud.CreateNetworkRulePortRequest, error) {
	var ports []ecloud.CreateNetworkRulePortRequest

	for _, rawPort := range rawPorts {
		port := rawPort.(map[string]interface{})
		protocol := port["protocol"].(string)
		protocolParsed, err := ecloud.ParseNetworkRulePortProtocol(protocol)
		if err != nil {
			return nil, err
		}

		ports = append(ports, ecloud.CreateNetworkRulePortRequest{
			Protocol:    protocolParsed,
			Source:      port["source"].(string),
			Destination: port["destination"].(string),
			Name:        port["name"].(string),
		})
	}

	return ports, nil
}

func expandUpdateNetworkRuleRequestPorts(rawPorts []interface{}) ([]ecloud.PatchNetworkRulePortRequest, error) {
	var ports []ecloud.PatchNetworkRulePortRequest

	for _, rawPort := range rawPorts {
		port := rawPort.(map[string]interface{})
		protocol := port["protocol"].(string)
		protocolParsed, err := ecloud.ParseNetworkRulePortProtocol(protocol)
		if err != nil {
			return nil, err
		}

		ports = append(ports, ecloud.PatchNetworkRulePortRequest{
			Protocol:    protocolParsed,
			Source:      port["source"].(string),
			Destination: port["destination"].(string),
			Name:        port["name"].(string),
		})
	}

	return ports, nil
}
