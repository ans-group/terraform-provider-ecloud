package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNATOverloadRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNATOverloadRuleRead,

		Schema: map[string]*schema.Schema{
			"nat_overload_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNATOverloadRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("nat_overload_rule_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if routerID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{routerID.(string)}))
	}
	if fipID, ok := d.GetOk("floating_ip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("floating_ip_id", connection.EQOperator, []string{fipID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	rules, err := service.GetNATOverloadRules(params)
	if err != nil {
		return fmt.Errorf("Error retrieving NAT overload rules: %s", err)
	}

	if len(rules) < 1 {
		return errors.New("No NAT overload rules found with provided arguments")
	}

	if len(rules) > 1 {
		return errors.New("More than 1 NAT overload rule found with provided arguments")
	}

	d.SetId(rules[0].ID)
	d.Set("network_id", rules[0].NetworkID)
	d.Set("floating_ip_id", rules[0].FloatingIPID)
	d.Set("subnet", rules[0].Subnet)
	d.Set("name", rules[0].Name)
	d.Set("action", rules[0].Action)

	return nil
}
