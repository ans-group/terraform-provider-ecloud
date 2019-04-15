package ecloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourcePod() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePodRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePodRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	pods, err := service.GetPods(params)
	if err != nil {
		return fmt.Errorf("Error retrieving pod: %s", err)
	}

	if len(pods) < 1 {
		return fmt.Errorf("Pod not found with name [%s]", name)
	}

	if len(pods) > 1 {
		return fmt.Errorf("More than one pod found with name [%s]", name)
	}

	d.SetId(strconv.Itoa(pods[0].ID))
	d.Set("name", pods[0].Name)

	return nil
}
