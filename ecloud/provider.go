package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/client"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/terraform-provider-ecloud/pkg/logger"
)

const userAgent = "terraform-provider-ecloud"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"context": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Config context to use",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API token to authenticate with UKFast APIs. See https://developers.ukfast.io for more details",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ecloud_volume":                    dataSourceVolume(),
			"ecloud_firewallpolicy":            dataSourceFirewallPolicy(),
			"ecloud_firewallrule":              dataSourceFirewallRule(),
			"ecloud_image":                     dataSourceImage(),
			"ecloud_instance":                  dataSourceInstance(),
			"ecloud_ipaddress":                 dataSourceIPAddress(),
			"ecloud_network":                   dataSourceNetwork(),
			"ecloud_router":                    dataSourceRouter(),
			"ecloud_router_throughput":         dataSourceRouterThroughput(),
			"ecloud_vpc":                       dataSourceVPC(),
			"ecloud_floatingip":                dataSourceFloatingIP(),
			"ecloud_nic":                       dataSourceNic(),
			"ecloud_hostspec":                  dataSourceHostSpec(),
			"ecloud_hostgroup":                 dataSourceHostGroup(),
			"ecloud_host":                      dataSourceHost(),
			"ecloud_ssh_keypair":               dataSourceSshKeyPair(),
			"ecloud_networkpolicy":             dataSourceNetworkPolicy(),
			"ecloud_networkrule":               dataSourceNetworkRule(),
			"ecloud_availability_zone":         dataSourceAvailabilityZone(),
			"ecloud_region":                    dataSourceRegion(),
			"ecloud_vpn_profile_group":         dataSourceVPNProfileGroup(),
			"ecloud_vpn_service":               dataSourceVPNService(),
			"ecloud_vpn_endpoint":              dataSourceVPNEndpoint(),
			"ecloud_vpn_session":               dataSourceVPNSession(),
			"ecloud_vpn_gateway":               dataSourceVPNGateway(),
			"ecloud_vpn_gateway_user":          dataSourceVPNGatewayUser(),
			"ecloud_vpn_gateway_specification": dataSourceVPNGatewaySpecification(),
			"ecloud_volumegroup":               dataSourceVolumeGroup(),
			"ecloud_loadbalancer_spec":         dataSourceLoadBalancerSpec(),
			"ecloud_loadbalancer":              dataSourceLoadBalancer(),
			"ecloud_loadbalancer_vip":          dataSourceLoadBalancerVip(),
			"ecloud_affinityrule":              dataSourceAffinityRule(),
			"ecloud_affinityrule_member":       dataSourceAffinityRuleMember(),
			"ecloud_resourcetier":              dataSourceResourceTier(),
			"ecloud_natoverloadrule":           dataSourceNATOverloadRule(),
			"ecloud_iops":                      dataSourceIOPS(),
			"ecloud_instance_credential":       dataSourceInstanceCredentials(),
			"ecloud_backup_gateway_spec":       dataSourceBackupGatewaySpecification(),
			"ecloud_backup_gateway":            dataSourceBackupGateway(),
			"ecloud_monitoring_gateway":        dataSourceMonitoringGateway(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ecloud_vpc":                   resourceVPC(),
			"ecloud_router":                resourceRouter(),
			"ecloud_network":               resourceNetwork(),
			"ecloud_image":                 resourceImage(),
			"ecloud_instance":              resourceInstance(),
			"ecloud_ipaddress":             resourceIPAddress(),
			"ecloud_firewallpolicy":        resourceFirewallPolicy(),
			"ecloud_firewallrule":          resourceFirewallRule(),
			"ecloud_volume":                resourceVolume(),
			"ecloud_floatingip":            resourceFloatingIP(),
			"ecloud_hostgroup":             resourceHostGroup(),
			"ecloud_host":                  resourceHost(),
			"ecloud_ssh_keypair":           resourceSshKeyPair(),
			"ecloud_networkpolicy":         resourceNetworkPolicy(),
			"ecloud_networkrule":           resourceNetworkRule(),
			"ecloud_nic_ipaddress_binding": resourceNICIPAddressBinding(),
			"ecloud_vpn_service":           resourceVPNService(),
			"ecloud_vpn_endpoint":          resourceVPNEndpoint(),
			"ecloud_vpn_session":           resourceVPNSession(),
			"ecloud_vpn_gateway":           resourceVPNGateway(),
			"ecloud_vpn_gateway_user":      resourceVPNGatewayUser(),
			"ecloud_volumegroup":           resourceVolumeGroup(),
			"ecloud_loadbalancer":          resourceLoadBalancer(),
			"ecloud_loadbalancer_vip":      resourceLoadBalancerVip(),
			"ecloud_affinityrule":          resourceAffinityRule(),
			"ecloud_affinityrule_member":   resourceAffinityRuleMember(),
			"ecloud_natoverloadrule":       resourceNATOverloadRule(),
			"ecloud_volumegroup_instance":  resourceVolumeGroupInstance(),
			"ecloud_instance_script":       resourceInstanceScript(),
			"ecloud_backup_gateway":        resourceBackupGateway(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	err := config.Init("")
	if err != nil {
		return nil, diag.Errorf("Failed to initialise config: %s", err)
	}

	if config.GetBool("api_debug") {
		logging.SetLogger(&logger.ProviderLogger{})
	}

	context := d.Get("context").(string)
	if len(context) > 0 {
		err := config.SwitchCurrentContext(context)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	apiKey := d.Get("api_key").(string)
	if len(apiKey) > 0 {
		config.Set(config.GetCurrentContextName(), "api_key", apiKey)
	}

	conn, err := getConnection()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client.NewClient(conn).ECloudService(), nil
}

func getConnection() (connection.Connection, error) {
	connFactory := connection.NewDefaultConnectionFactory(
		connection.WithDefaultConnectionUserAgent(userAgent),
	)

	return connFactory.NewConnection()
}
