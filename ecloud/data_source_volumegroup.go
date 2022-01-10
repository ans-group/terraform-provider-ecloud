package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceVolumeGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVolumeGroupRead,

		Schema: map[string]*schema.Schema{
			"volume_group_id": {
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
		},
	}
}

func dataSourceVolumeGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("volume_group_id"); ok {
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

	volumegroups, err := service.GetVolumeGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active volumesgroups: %s", err)
	}

	if len(volumegroups) < 1 {
		return errors.New("No volumesgroups found with provided arguments")
	}

	if len(volumegroups) > 1 {
		return errors.New("More than 1 volumegroups found with provided arguments")
	}

	d.SetId(volumegroups[0].ID)
	d.Set("name", volumegroups[0].Name)
	d.Set("vpc_id", volumegroups[0].VPCID)
	d.Set("availability_zone_id", volumegroups[0].AvailabilityZoneID)

	return nil
}
