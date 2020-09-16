package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallRuleRead,

		Schema: map[string]*schema.Schema{
			"firewallrule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("firewallrule_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if routerID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{routerID.(string)}))
	}

	azs, err := service.GetFirewallRules(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active firewall rules: %s", err)
	}

	if len(azs) < 1 {
		return errors.New("No firewall rules found with provided arguments")
	}

	if len(azs) > 1 {
		return errors.New("More than 1 firewall rule found with provided arguments")
	}

	d.SetId(azs[0].ID)
	d.Set("router_id", azs[0].RouterID)

	return nil
}
