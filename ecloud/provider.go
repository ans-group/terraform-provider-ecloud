package ecloud

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: func() (interface{}, error) {
					key := os.Getenv("UKF_API_KEY")
					if key != "" {
						return key, nil
					}

					return "", errors.New("api_key required")
				},
				Description: "API token required to authenticate with UKFast APIs. See https://developers.ukfast.io for more details",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ecloud_datastore":         dataSourceDatastore(),
			"ecloud_site":              dataSourceSite(),
			"ecloud_solution":          dataSourceSolution(),
			"ecloud_network":           dataSourceNetwork(),
			"ecloud_pod":               dataSourcePod(),
			"ecloud_template":          dataSourceTemplate(),
			"ecloud_solution_template": dataSourceSolutionTemplate(),
			"ecloud_appliance":         dataSourceAppliance(),
			"ecloud_pod_appliance":     dataSourcePodAppliance(),
			"ecloud_active_directory":  dataSourceActiveDirectory(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ecloud_virtualmachine":     resourceVirtualMachine(),
			"ecloud_virtualmachine_tag": resourceVirtualMachineTag(),
			"ecloud_solution_tag":       resourceSolutionTag(),
			"ecloud_solution_template":  resourceSolutionTemplate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return getService(d.Get("api_key").(string)), nil
}

func getClient(apiKey string) client.Client {
	return client.NewClient(connection.NewAPIKeyCredentialsAPIConnection(apiKey))
}

func getService(apiKey string) ecloudservice.ECloudService {
	return getClient(apiKey).ECloudService()
}
