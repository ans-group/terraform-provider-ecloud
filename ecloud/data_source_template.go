package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hdd": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)
	podID := d.Get("pod_id").(int)

	templates, err := service.GetPodTemplates(podID, connection.APIRequestParameters{})
	if err != nil {
		return fmt.Errorf("Error retrieving pod templates: %s", err)
	}

	template, err := filterTemplateName(templates, name)
	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("name", name)
	d.Set("pod_id", podID)
	d.Set("cpu", template.CPU)
	d.Set("ram", template.RAM)
	d.Set("hdd", template.HDD)
	d.Set("platform", template.Platform)

	return nil
}

func filterTemplateName(templates []ecloudservice.Template, name string) (ecloudservice.Template, error) {
	var foundTemplates []ecloudservice.Template
	for _, template := range templates {
		if template.Name == name {
			foundTemplates = append(foundTemplates, template)
		}
	}

	if len(foundTemplates) < 1 {
		return ecloudservice.Template{}, fmt.Errorf("Template not found with name [%s]", name)
	}

	if len(foundTemplates) > 1 {
		return ecloudservice.Template{}, fmt.Errorf("More than one template found with name [%s]", name)
	}

	return foundTemplates[0], nil
}
