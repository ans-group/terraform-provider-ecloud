package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRouterThroughput() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouterThroughputRead,

		Schema: map[string]*schema.Schema{
			"router_throughput_id": {
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
			"committed_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRouterThroughputRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("router_throughput_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	throughputs, err := service.GetRouterThroughputs(params)
	if err != nil {
		return diag.Errorf("Error retrieving router throughputs: %s", err)
	}

	if len(throughputs) < 1 {
		return diag.Errorf("No router throughputs found with provided arguments")
	}

	if len(throughputs) > 1 {
		return diag.Errorf("More than 1 router throughput found with provided arguments")
	}

	d.SetId(throughputs[0].ID)
	d.Set("availability_zone_id", throughputs[0].AvailabilityZoneID)
	d.Set("name", throughputs[0].Name)
	d.Set("committed_bandwidth", throughputs[0].CommittedBandwidth)

	return nil
}
