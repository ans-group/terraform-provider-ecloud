package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNSession() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNSessionRead,

		Schema: map[string]*schema.Schema{
			"vpn_session_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_profile_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_networks": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_networks": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"psk": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNSessionRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_session_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpnServiceID, ok := d.GetOk("vpn_service_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpn_service_id", connection.EQOperator, []string{vpnServiceID.(string)}))
	}
	if vpnProfileGroupID, ok := d.GetOk("vpn_profile_group_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpn_profile_group_id", connection.EQOperator, []string{vpnProfileGroupID.(string)}))
	}
	if vpnEndpointID, ok := d.GetOk("vpn_endpoint_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpn_endpoint_id", connection.EQOperator, []string{vpnEndpointID.(string)}))
	}
	if remoteIP, ok := d.GetOk("remote_ip"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("remote_ip", connection.EQOperator, []string{remoteIP.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	vpnSessions, err := service.GetVPNSessions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPN services: %s", err)
	}

	if len(vpnSessions) < 1 {
		return errors.New("No VPN services found with provided arguments")
	}

	if len(vpnSessions) > 1 {
		return errors.New("More than 1 VPN service found with provided arguments")
	}

	d.SetId(vpnSessions[0].ID)
	d.Set("vpn_service_id", vpnSessions[0].VPNServiceID)
	d.Set("vpn_profile_group_id", vpnSessions[0].VPNProfileGroupID)
	d.Set("vpn_endpoint_id", vpnSessions[0].VPNEndpointID)
	d.Set("remote_ip", vpnSessions[0].RemoteIP)
	d.Set("name", vpnSessions[0].Name)
	d.Set("remote_networks", vpnSessions[0].RemoteNetworks)
	d.Set("local_networks", vpnSessions[0].LocalNetworks)

	psk, err := service.GetVPNSessionPreSharedKey(vpnSessions[0].ID)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN service pre-shared key: %s", err)
	}
	d.Set("psk", psk.PSK)

	return nil
}
