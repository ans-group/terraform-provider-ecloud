package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceSolutionTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSolutionTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"solution_id": &schema.Schema{
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

func dataSourceSolutionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)
	solutionID := d.Get("solution_id").(int)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	templates, err := service.GetSolutionTemplates(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution templates: %s", err)
	}

	template, err := filterTemplateName(templates, name)
	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("solution_id", solutionID)
	d.Set("name", name)
	d.Set("cpu", template.CPU)
	d.Set("ram", template.RAM)
	d.Set("hdd", template.HDD)
	d.Set("platform", template.Platform)

	return nil
}
