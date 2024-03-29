package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFloatingIPRead,

		Schema: map[string]*schema.Schema{
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone_id": {
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
		},
	}
}

func dataSourceFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("floating_ip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if ip, ok := d.GetOk("ip_address"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("ip_address", connection.EQOperator, []string{ip.(string)}))
	}

	fips, err := service.GetFloatingIPs(params)
	if err != nil {
		return diag.Errorf("Error retrieving floating ips: %s", err)
	}

	if len(fips) < 1 {
		return diag.Errorf("No floating ips found with provided arguments")
	}

	if len(fips) > 1 {
		return diag.Errorf("More than 1 floating ip found with provided arguments")
	}

	d.SetId(fips[0].ID)
	d.Set("vpc_id", fips[0].VPCID)
	d.Set("availability_zone_id", fips[0].AvailabilityZoneID)
	d.Set("name", fips[0].Name)
	d.Set("ip_address", fips[0].IPAddress)

	return nil
}
