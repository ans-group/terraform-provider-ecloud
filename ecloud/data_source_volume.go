package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVolumeRead,

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("volume_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if capacity, ok := d.GetOk("capacity"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("capacity", connection.EQOperator, []string{strconv.Itoa(capacity.(int))}))
	}
	if iops, ok := d.GetOk("iops"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("iops", connection.EQOperator, []string{strconv.Itoa(iops.(int))}))
	}

	volumes, err := service.GetVolumes(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active volumes: %s", err)
	}

	if len(volumes) < 1 {
		return errors.New("No volumes found with provided arguments")
	}

	if len(volumes) > 1 {
		return errors.New("More than 1 image found with provided arguments")
	}

	d.SetId(volumes[0].ID)
	d.Set("name", volumes[0].Name)
	d.Set("capacity", volumes[0].Capacity)
	d.Set("iops", volumes[0].IOPS)

	return nil
}
