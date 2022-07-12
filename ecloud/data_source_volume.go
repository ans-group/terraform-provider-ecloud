package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVolumeRead,

		Schema: map[string]*schema.Schema{
			"volume_id": {
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
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"volume_group_id": {
				Type:     schema.TypeString,
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
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
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
	if volumeGroupID, ok := d.GetOk("volume_group_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("volume_group_id", connection.EQOperator, []string{volumeGroupID.(string)}))
	}

	volumes, err := service.GetVolumes(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active volumes: %s", err)
	}

	if len(volumes) < 1 {
		return errors.New("No volumes found with provided arguments")
	}

	if len(volumes) > 1 {
		return errors.New("More than 1 volume found with provided arguments")
	}

	d.SetId(volumes[0].ID)
	d.Set("name", volumes[0].Name)
	d.Set("capacity", volumes[0].Capacity)
	d.Set("iops", volumes[0].IOPS)
	d.Set("vpc_id", volumes[0].VPCID)
	d.Set("availability_zone_id", volumes[0].AvailabilityZoneID)
	d.Set("volume_group_id", volumes[0].VolumeGroupID)

	return nil
}
