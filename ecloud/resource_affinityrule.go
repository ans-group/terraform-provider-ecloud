package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAffinityRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAffinityRuleCreate,
		ReadContext:   resourceAffinityRuleRead,
		UpdateContext: resourceAffinityRuleUpdate,
		DeleteContext: resourceAffinityRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := ecloudservice.ParseAffinityRuleType(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid affinity rule type [affinity, anti-affinity], got: %s", key, v))
					}
					return
				},
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceAffinityRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	ruleType, err := ecloudservice.ParseAffinityRuleType(d.Get("type").(string))
	if err != nil {
		return diag.Errorf("Error parsing affinity rule type: %s", err)
	}

	createReq := ecloudservice.CreateAffinityRuleRequest{
		VPCID:              d.Get("vpc_id").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
		Name:               d.Get("name").(string),
		Type:               ruleType,
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateAffinityRuleRequest: %+v", createReq))

	tflog.Info(ctx, "Creating AffinityRule")
	taskRef, err := service.CreateAffinityRule(createReq)
	if err != nil {
		return diag.Errorf("Error creating affinity rule: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for affinity rule with ID [%s] to be created: %s", d.Id(), err)
	}

	// add members to rule
	if rawIDs, ok := d.GetOk("instance_ids"); ok {
		for _, rawID := range rawIDs.(*schema.Set).List() {
			memberID := rawID.(string)
			if len(memberID) < 1 {
				continue
			}

			req := ecloudservice.CreateAffinityRuleMemberRequest{
				AffinityRuleID: d.Id(),
				InstanceID:     memberID,
			}

			tflog.Debug(ctx, fmt.Sprintf("Created CreateAffinityRuleMemberRequest: %+v", req))

			tflog.Info(ctx, "Adding instance member to affinity rule", map[string]interface{}{
				"affinity_rule_id": d.Id(),
				"instance_member":  memberID,
			})
			taskRef, err := service.CreateAffinityRuleMember(req)
			if err != nil {
				return diag.Errorf("Error creating affinity rule member: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for affinity rule member with ID [%s] to be created: %s", d.Id(), err)
			}
		}
	}

	return resourceAffinityRuleRead(ctx, d, meta)
}

func resourceAffinityRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving affinity rule", map[string]interface{}{
		"id": d.Id(),
	})
	ar, err := service.GetAffinityRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	ruleMembers, err := service.GetAffinityRuleMembers(d.Id(), connection.APIRequestParameters{})
	if err != nil {
		return diag.Errorf("Error retrieving affinity rule members for rule ID [%s]: %s", d.Id(), err)
	}

	d.Set("vpc_id", ar.VPCID)
	d.Set("name", ar.Name)
	d.Set("availability_zone_id", ar.AvailabilityZoneID)
	d.Set("type", ar.Type)
	d.Set("instance_ids", flattenAffinityRuleMembers(ruleMembers))

	return nil
}

func resourceAffinityRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating affinity rule", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchAffinityRuleRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchAffinityRule(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating affinity rule with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for affinity rule with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("instance_ids") {
		oldRaw, newRaw := d.GetChange("instance_ids")

		oldIDs := oldRaw.(*schema.Set).List()
		newIDs := newRaw.(*schema.Set).List()

		ruleMembers, err := service.GetAffinityRuleMembers(d.Id(), connection.APIRequestParameters{})
		if err != nil {
			return diag.Errorf("Error retrieving affinity rule members for rule ID [%s]: %s", d.Id(), err)
		}
		getMemberByInstanceID := func(slice []ecloudservice.AffinityRuleMember, value string) string {
			for _, s := range slice {
				if s.InstanceID == value {
					return s.ID
				}
			}
			return ""
		}

		for _, id := range oldIDs {
			instanceID := id.(string)
			if rawMemberExistsById(newIDs, instanceID) {
				continue
			}

			ruleMemberID := getMemberByInstanceID(ruleMembers, instanceID)

			tflog.Info(ctx, "Removing instance member to affinity rule", map[string]interface{}{
				"affinity_rule_id": d.Id(),
				"instance_member":  instanceID,
			})

			taskID, err := service.DeleteAffinityRuleMember(ruleMemberID)
			if err != nil {
				return diag.Errorf("Error deleting affinity rule member: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for affinity rule member for instance ID [%s] to be delete: %s", instanceID, err)
			}
		}

		for _, id := range newIDs {
			instanceID := id.(string)
			if rawMemberExistsById(oldIDs, instanceID) {
				continue
			}

			req := ecloudservice.CreateAffinityRuleMemberRequest{
				AffinityRuleID: d.Id(),
				InstanceID:     instanceID,
			}

			tflog.Debug(ctx, fmt.Sprintf("Created CreateAffinityRuleMemberRequest: %+v", req))

			tflog.Info(ctx, "Adding instance member to affinity rule", map[string]interface{}{
				"affinity_rule_id": d.Id(),
				"instance_member":  instanceID,
			})
			taskRef, err := service.CreateAffinityRuleMember(req)
			if err != nil {
				return diag.Errorf("Error creating affinity rule member: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for affinity rule member for instance ID [%s] to be created: %s", instanceID, err)
			}

		}
	}

	return resourceAffinityRuleRead(ctx, d, meta)
}

func resourceAffinityRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	// delete members first if they exist
	if memberInstances, ok := d.GetOk("instance_ids"); ok {
		for _, memberInstance := range memberInstances.(*schema.Set).List() {
			instanceID := memberInstance.(string)

			if len(instanceID) < 1 {
				continue
			}

			// get affinity rule member ID by instanceID
			params := connection.APIRequestParameters{}
			params.WithFilter(*connection.NewAPIRequestFiltering("instance_id", connection.EQOperator, []string{instanceID}))

			arMembers, err := service.GetAffinityRuleMembers(d.Id(), params)
			if err != nil {
				return diag.Errorf("Error retrieving affinity rule members for rule ID [%s]: %s", d.Id(), err)
			}

			if len(arMembers) < 1 {
				// resource may have been deleted already, so skip
				continue
			}

			if len(arMembers) > 1 {
				return diag.Errorf("More than 1 affinity rule member found for instance ID [%s]", instanceID)
			}

			tflog.Info(ctx, "Adding instance member to affinity rule", map[string]interface{}{
				"affinity_rule_id": arMembers[0].ID,
				"instance_member":  instanceID,
			})
			taskID, err := service.DeleteAffinityRuleMember(arMembers[0].ID)
			if err != nil {
				switch err.(type) {
				case *ecloudservice.AffinityRuleMemberNotFoundError:
					tflog.Debug(ctx, "Affinity rule member not found, skipping delete", map[string]interface{}{
						"affinity_rule_id": arMembers[0].ID,
						"instance_member":  instanceID,
					})
				default:
					return diag.Errorf("Error removing affinity rule member ID [%s] for instance ID [%s]: %s", arMembers[0].ID, instanceID, err)
				}
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for affinity rule member ID [%s] for instance ID [%s] to be removed: %s", arMembers[0].ID, instanceID, err)
			}
		}
	}

	tflog.Info(ctx, "Removing affinity rule", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteAffinityRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing affinity rule with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for affinity rule with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
