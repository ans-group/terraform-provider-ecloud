package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIOPS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIOPSRead,

		Schema: map[string]*schema.Schema{
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func dataSourceIOPSRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("iops_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if level, ok := d.GetOk("level"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("level", connection.EQOperator, []string{strconv.Itoa(level.(int))}))
	}

	var tiers []ecloudservice.IOPSTier
	var err error
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		tiers, err = service.GetAvailabilityZoneIOPSTiers(azID.(string), params)
		if err != nil {
			return fmt.Errorf("Error retrieving availability zone IOPS tiers: %s", err)
		}
	} else {
		tiers, err = service.GetIOPSTiers(params)
		if err != nil {
			return fmt.Errorf("Error retrieving IOPS tiers: %s", err)
		}
	}

	if len(tiers) < 1 {
		return errors.New("No IOPS tiers found with provided arguments")
	}

	if len(tiers) > 1 {
		return errors.New("More than 1 IOPS tier found with provided arguments")
	}

	d.SetId(tiers[0].ID)
	d.Set("name", tiers[0].Name)
	d.Set("level", tiers[0].Level)

	return nil
}
