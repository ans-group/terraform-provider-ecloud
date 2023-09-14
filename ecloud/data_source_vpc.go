package ecloud

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPCRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if regionID, ok := d.GetOk("region_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("region_id", connection.EQOperator, []string{regionID.(string)}))
	}
	if clientID, ok := d.GetOk("client_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("client_id", connection.EQOperator, []string{strconv.Itoa(clientID.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	vpcs, err := service.GetVPCs(params)
	if err != nil {
		return diag.Errorf("Error retrieving active VPCs: %s", err)
	}

	if len(vpcs) < 1 {
		return diag.Errorf("No VPCs found with provided arguments")
	}

	if len(vpcs) > 1 {
		return diag.Errorf("More than 1 VPC found with provided arguments")
	}

	d.SetId(vpcs[0].ID)
	d.Set("region_id", vpcs[0].RegionID)
	d.Set("client_id", vpcs[0].ClientID)
	d.Set("name", vpcs[0].Name)

	return nil
}
