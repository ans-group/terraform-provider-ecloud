package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceRouterThroughput() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRouterThroughputRead,

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
		},
	}
}

func dataSourceRouterThroughputRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving router throughputs: %s", err)
	}

	if len(throughputs) < 1 {
		return errors.New("No router throughputs found with provided arguments")
	}

	if len(throughputs) > 1 {
		return errors.New("More than 1 router throughput found with provided arguments")
	}

	d.SetId(throughputs[0].ID)
	d.Set("availability_zone_id", throughputs[0].AvailabilityZoneID)
	d.Set("name", throughputs[0].Name)
	d.Set("committed_bandwidth", throughputs[0].CommittedBandwidth)
	d.Set("burst_size", throughputs[0].BurstSize)

	return nil
}
