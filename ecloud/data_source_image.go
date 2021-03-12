package ecloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"image_id": {
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

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("image_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}

	images, err := service.GetImages(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active images: %s", err)
	}

	if name, ok := d.GetOk("name"); ok {
		images = filterImageName(images, name.(string))
		if err != nil {
			return err
		}
	}

	if len(images) < 1 {
		return errors.New("No images found with provided arguments")
	}

	if len(images) > 1 {
		return errors.New("More than 1 image found with provided arguments")
	}

	d.SetId(images[0].ID)
	d.Set("name", images[0].Name)

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
