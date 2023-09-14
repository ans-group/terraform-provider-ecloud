package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAffinityRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAffinityRuleRead,

		Schema: map[string]*schema.Schema{
			"affinity_rule_id": {
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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAffinityRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("affinity_rule_id"); ok {
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
	if ruleType, ok := d.GetOk("type"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("type", connection.EQOperator, []string{ruleType.(string)}))
	}

	ars, err := service.GetAffinityRules(params)
	if err != nil {
		return diag.Errorf("Error retrieving affinity rules: %s", err)
	}

	if len(ars) < 1 {
		return diag.Errorf("No affinity rules found with provided arguments")
	}

	if len(ars) > 1 {
		return diag.Errorf("More than 1 affinity rule found with provided arguments")
	}

	d.SetId(ars[0].ID)
	d.Set("vpc_id", ars[0].VPCID)
	d.Set("availability_zone_id", ars[0].AvailabilityZoneID)
	d.Set("type", ars[0].Type)
	d.Set("name", ars[0].Name)

	return nil
}
