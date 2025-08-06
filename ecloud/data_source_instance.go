package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("instance_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	instances, err := service.GetInstances(params)
	if err != nil {
		return diag.Errorf("Error retrieving active instances: %s", err)
	}

	if len(instances) < 1 {
		return diag.Errorf("No instances found with provided arguments")
	}

	if len(instances) > 1 {
		return diag.Errorf("More than 1 instance found with provided arguments")
	}

	d.SetId(instances[0].ID)
	d.Set("vpc_id", instances[0].VPCID)
	d.Set("name", instances[0].Name)

	// Set tags
	if err := d.Set("tags", flattenInstanceTags(instances[0].Tags)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
