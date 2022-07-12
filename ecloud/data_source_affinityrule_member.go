package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAffinityRuleMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAffinityRuleMemberRead,

		Schema: map[string]*schema.Schema{
			"affinity_rule_member_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"affinity_rule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAffinityRuleMemberRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if memberID, ok := d.GetOk("affinity_rule_member_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{memberID.(string)}))
	}
	if ruleID, ok := d.GetOk("affinity_rule_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("affinity_rule_id", connection.EQOperator, []string{ruleID.(string)}))
	}
	if instanceID, ok := d.GetOk("instance_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("instance_id", connection.EQOperator, []string{instanceID.(string)}))
	}

	arMembers, err := service.GetAffinityRuleMembers(d.Get("affinity_rule_id").(string), params)
	if err != nil {
		return fmt.Errorf("Error retrieving affinity rule members: %s", err)
	}

	if len(arMembers) < 1 {
		return errors.New("No affinity rule members found with provided arguments")
	}

	if len(arMembers) > 1 {
		return errors.New("More than 1 affinity rule member found with provided arguments")
	}

	d.SetId(arMembers[0].ID)
	d.Set("affinity_rule_id", arMembers[0].AffinityRuleID)
	d.Set("instance_id", arMembers[0].InstanceID)

	return nil
}
