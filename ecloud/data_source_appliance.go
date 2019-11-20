package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceAppliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplianceRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"appliance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceApplianceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	appliances, err := service.GetAppliances(connection.APIRequestParameters{})

	if err != nil {
		return fmt.Errorf("Error retrieving appliances: %s", err)
	}

	appliance, err := filterApplianceName(appliances, name)
	if err != nil {
		return err
	}

	d.SetId(appliance.ID)
	d.Set("name", appliance.Name)
	d.Set("appliance_id", appliance.ID)

	return nil
}

func filterApplianceName(appliances []ecloudservice.Appliance, name string) (ecloudservice.Appliance, error) {
	var foundAppliances []ecloudservice.Appliance
	for _, appliance := range appliances {
		if appliance.Name == name {
			foundAppliances = append(foundAppliances, appliance)
		}
	}

	if len(foundAppliances) < 1 {
		return ecloudservice.Appliance{}, fmt.Errorf("Appliance not found with name [%s]", name)
	}
	if len(foundAppliances) > 1 {
		return ecloudservice.Appliance{}, fmt.Errorf("More than one appliance found with name [%s]", name)
	}

	return foundAppliances[0], nil
}
