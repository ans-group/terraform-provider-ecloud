package ecloud

import (
	"context"
	"strings"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
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
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("image_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}

	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}

	if availabilityZoneID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{availabilityZoneID.(string)}))
	}

	if platform, ok := d.GetOk("platform"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("platform", connection.EQOperator, []string{platform.(string)}))
	}

	images, err := service.GetImages(params)
	if err != nil {
		return diag.Errorf("Error retrieving active images: %s", err)
	}

	if name, ok := d.GetOk("name"); ok {
		images = filterImageName(images, name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(images) < 1 {
		return diag.Errorf("No images found with provided arguments")
	}

	if len(images) > 1 {
		return diag.Errorf("More than 1 image found with provided arguments")
	}

	d.SetId(images[0].ID)
	d.Set("name", images[0].Name)
	d.Set("vpc_id", images[0].VPCID)
	d.Set("availability_zone_id", images[0].AvailabilityZoneID)
	d.Set("platform", images[0].Platform)

	return nil
}

func filterImageName(images []ecloudservice.Image, name string) []ecloudservice.Image {
	for _, image := range images {
		if strings.ToLower(image.Name) == strings.ToLower(name) {
			return []ecloudservice.Image{image}
		}
	}

	return []ecloudservice.Image{}
}
