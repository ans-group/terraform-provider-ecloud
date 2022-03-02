package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceHostSpec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostSpecRead,

		Schema: map[string]*schema.Schema{
			"host_spec_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceHostSpecRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("host_spec_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	hostSpecs, err := service.GetHostSpecs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving host specs: %s", err)
	}

	if len(hostSpecs) < 1 {
		return errors.New("No host specs found with provided arguments")
	}

	if len(hostSpecs) > 1 {
		return errors.New("More than 1 host spec found with provided arguments")
	}

	d.SetId(hostSpecs[0].ID)
	d.Set("name", hostSpecs[0].Name)
	d.Set("cpu_sockets", hostSpecs[0].CPUSockets)
	d.Set("cpu_cores", hostSpecs[0].CPUCores)
	d.Set("ram_capacity", hostSpecs[0].RAMCapacity)

	return nil
}
