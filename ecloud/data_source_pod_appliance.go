package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourcePodAppliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePodApplianceRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeInt,
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

func dataSourcePodApplianceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)
	podID := d.Get("pod_id").(int)

	appliances, err := service.GetPodAppliances(podID, connection.APIRequestParameters{})

	if err != nil {
		return fmt.Errorf("Error retrieving pod appliances: %s", err)
	}

	appliance, err := filterApplianceName(appliances, name)
	if err != nil {
		return err
	}

	d.SetId(appliance.ID)
	d.Set("name", appliance.Name)
	d.Set("appliance_id", appliance.ID)
	d.Set("pod_id", podID)

	return nil
}
