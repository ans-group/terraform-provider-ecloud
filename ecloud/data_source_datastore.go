package ecloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceDatastore() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatastoreRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"solution_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)
	solutionID := d.Get("solution_id").(int)
	siteID := d.Get("site_id").(int)

	params := connection.APIRequestParameters{}

	if solutionID > 0 {
		params.WithFilter(connection.APIRequestFiltering{
			Property: "solution_id",
			Operator: connection.EQOperator,
			Value:    []string{strconv.Itoa(solutionID)},
		})
	}

	if siteID > 0 {
		params.WithFilter(connection.APIRequestFiltering{
			Property: "site_id",
			Operator: connection.EQOperator,
			Value:    []string{strconv.Itoa(siteID)},
		})
	}

	datastores, err := service.GetDatastores(params)
	if err != nil {
		return fmt.Errorf("Error retrieving datastores: %s", err)
	}

	datastore, err := filterDatastoreName(datastores, name)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(datastore.ID))
	d.Set("name", datastore.Name)
	d.Set("solution_id", datastore.SolutionID)
	d.Set("site_id", datastore.SiteID)
	d.Set("status", datastore.Status)
	d.Set("capacity", datastore.Capacity)

	return nil
}

func filterDatastoreName(datastores []ecloudservice.Datastore, name string) (ecloudservice.Datastore, error) {
	var foundDatastores []ecloudservice.Datastore
	for _, datastore := range datastores {
		if datastore.Name == name {
			foundDatastores = append(foundDatastores, datastore)
		}
	}

	if len(foundDatastores) < 1 {
		return ecloudservice.Datastore{}, fmt.Errorf("Datastore not found with name [%s]", name)
	}

	if len(foundDatastores) > 1 {
		return ecloudservice.Datastore{}, fmt.Errorf("More than one datastore found with name [%s]", name)
	}

	return foundDatastores[0], nil
}
