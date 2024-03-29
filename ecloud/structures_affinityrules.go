package ecloud

import (
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//flattenAffinityRuleMembers flattens affinity rule members into a set
func flattenAffinityRuleMembers(ruleMembers []ecloudservice.AffinityRuleMember) *schema.Set {
	memberIDs := schema.NewSet(schema.HashString, []interface{}{})
	for _, member := range ruleMembers {
		memberIDs.Add(member.InstanceID)
	}

	return memberIDs
}

//rawMemberExistsById returns true if value is in slice
func rawMemberExistsById(rawMembers []interface{}, value string) bool {
	for _, rawMember := range rawMembers {
		member := rawMember.(string)
		if member == value {
			return true
		}
	}

	return false
}
