package ecloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPAddressRead,

		Schema: map[string]*schema.Schema{
			"ip_address_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("ip_address_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if ipAddress, ok := d.GetOk("ip_address"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("ip_address", connection.EQOperator, []string{ipAddress.(string)}))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{networkID.(string)}))
	}
	if ipType, ok := d.GetOk("type"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("type", connection.EQOperator, []string{ipType.(string)}))
	}

	log.Printf("[DEBUG] Retrieving IP addresses with parameters: %+v", params)
	addresses, err := service.GetIPAddresses(params)
	if err != nil {
		return fmt.Errorf("Error retrieving IP addresses: %s", err)
	}

	if len(addresses) != 1 {
		return fmt.Errorf("Unexpected number [%d] of IP addresses found, expected 1", len(addresses))
	}

	d.SetId(addresses[0].ID)
	d.Set("name", addresses[0].Name)
	d.Set("ip_address", addresses[0].IPAddress)
	d.Set("network_id", addresses[0].NetworkID)
	d.Set("type", addresses[0].Type)

	return nil
}
