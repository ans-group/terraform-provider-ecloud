package ecloud

import (
	"errors"
	"fmt"
	"log"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkPolicyRead,

		Schema: map[string]*schema.Schema{
			"network_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNetworkPolicyRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("network_policy_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{networkID.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	log.Printf("[DEBUG] Retrieving network policies with parameters: %+v", params)
	policies, err := service.GetNetworkPolicies(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active network policies: %s", err)
	}

	if len(policies) < 1 {
		return errors.New("No network policies found with provided arguments")
	}

	if len(policies) > 1 {
		return errors.New("More than 1 network policy found with provided arguments")
	}

	d.SetId(policies[0].ID)
	d.Set("network_id", policies[0].NetworkID)
	d.Set("vpc_id", policies[0].VPCID)
	d.Set("name", policies[0].Name)

	return nil
}
