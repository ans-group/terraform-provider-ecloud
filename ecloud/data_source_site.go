package ecloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSiteRead,

		Schema: map[string]*schema.Schema{
			"pod_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"solution_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSiteRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	podID := d.Get("pod_id").(int)
	solutionID := d.Get("solution_id").(int)

	params := connection.APIRequestParameters{}

	params.WithFilter(connection.APIRequestFiltering{
		Property: "pod_id",
		Operator: connection.EQOperator,
		Value:    []string{strconv.Itoa(podID)},
	})

	params.WithFilter(connection.APIRequestFiltering{
		Property: "solution_id",
		Operator: connection.EQOperator,
		Value:    []string{strconv.Itoa(solutionID)},
	})

	sites, err := service.GetSites(params)
	if err != nil {
		return fmt.Errorf("Error retrieving sites: %s", err)
	}

	if len(sites) < 1 {
		return fmt.Errorf("No site found with specified pod_id [%d] and solution_id [%d]", podID, solutionID)
	}

	if len(sites) > 1 {
		return fmt.Errorf("More than one site found with specified pod_id [%d] and solution_id [%d]", podID, solutionID)
	}

	d.SetId(strconv.Itoa(sites[0].ID))
	d.Set("pod_id", sites[0].PodID)
	d.Set("solution_id", sites[0].SolutionID)
	d.Set("state", sites[0].State)

	return nil
}
