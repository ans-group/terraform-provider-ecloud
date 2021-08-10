package ecloud

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: func() (interface{}, error) {
					key := os.Getenv("UKF_API_KEY")
					if key != "" {
						return key, nil
					}

					return "", errors.New("api_key required")
				},
				Description: "API token required to authenticate with UKFast APIs. See https://developers.ukfast.io for more details",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ecloud_volume":            dataSourceVolume(),
			"ecloud_firewallpolicy":    dataSourceFirewallPolicy(),
			"ecloud_firewallrule":      dataSourceFirewallRule(),
			"ecloud_image":             dataSourceImage(),
			"ecloud_instance":          dataSourceInstance(),
			"ecloud_network":           dataSourceNetwork(),
			"ecloud_router":            dataSourceRouter(),
			"ecloud_router_throughput": dataSourceRouterThroughput(),
			"ecloud_vpc":               dataSourceVPC(),
			"ecloud_floatingip":        dataSourceFloatingIP(),
			"ecloud_nic":               dataSourceNic(),
			"ecloud_hostspec":          dataSourceHostSpec(),
			"ecloud_hostgroup":         dataSourceHostGroup(),
			"ecloud_host":              dataSourceHost(),
			"ecloud_ssh_keypair":       dataSourceSshKeyPair(),
			"ecloud_networkpolicy":     dataSourceNetworkPolicy(),
			"ecloud_networkrule":       dataSourceNetworkRule(),
			"ecloud_availability_zone": dataSourceAvailabilityZone(),
			"ecloud_region":            dataSourceRegion(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ecloud_vpc":            resourceVPC(),
			"ecloud_router":         resourceRouter(),
			"ecloud_network":        resourceNetwork(),
			"ecloud_instance":       resourceInstance(),
			"ecloud_firewallpolicy": resourceFirewallPolicy(),
			"ecloud_firewallrule":   resourceFirewallRule(),
			"ecloud_volume":         resourceVolume(),
			"ecloud_floatingip":     resourceFloatingIP(),
			"ecloud_hostgroup":      resourceHostGroup(),
			"ecloud_host":           resourceHost(),
			"ecloud_ssh_keypair":    resourceSshKeyPair(),
			"ecloud_networkpolicy":  resourceNetworkPolicy(),
			"ecloud_networkrule":    resourceNetworkRule(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return getService(d.Get("api_key").(string)), nil
}

func getClient(apiKey string) client.Client {
	return client.NewClient(connection.NewAPIKeyCredentialsAPIConnection(apiKey))
}

func getService(apiKey string) ecloudservice.ECloudService {
	return getClient(apiKey).ECloudService()
}
