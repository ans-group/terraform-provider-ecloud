package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceActiveDirectory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceActiveDirectoryRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ad_domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	domains, err := service.GetActiveDirectoryDomains(connection.APIRequestParameters{})

	if err != nil {
		return fmt.Errorf("Error retrieving active directory domains: %s", err)
	}

	domain, err := filterActiveDirectoryName(domains, name)
	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("name", domain.Name)
	d.Set("ad_domain_id", domain.ID)

	return nil
}

func filterActiveDirectoryName(domains []ecloudservice.ActiveDirectoryDomain, name string) (ecloudservice.ActiveDirectoryDomain, error) {
	var foundDomains []ecloudservice.ActiveDirectoryDomain
	for _, domain := range domains {
		if domain.Name == name {
			foundDomains = append(foundDomains, domain)
		}
	}

	if len(foundDomains) < 1 {
		return ecloudservice.ActiveDirectoryDomain{}, fmt.Errorf("Active Directory domain not found with name [%s]", name)
	}
	if len(foundDomains) > 1 {
		return ecloudservice.ActiveDirectoryDomain{}, fmt.Errorf("More than one Active Directory domain found with name [%s]", name)
	}

	return foundDomains[0], nil
}
