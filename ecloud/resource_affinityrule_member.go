package ecloud

import (
	"context"
	"fmt"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/terraform-provider-ecloud/pkg/lock"
)

func resourceAffinityRuleMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAffinityRuleMemberCreate,
		ReadContext:   resourceAffinityRuleMemberRead,
		UpdateContext: resourceAffinityRuleMemberUpdate,
		DeleteContext: resourceAffinityRuleMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"affinity_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAffinityRuleMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ruleID := d.Get("affinity_rule_id").(string)
	if len(ruleID) < 1 {
		return diag.Errorf("Invalid affinity rule ID: %s", ruleID)
	}

	unlock := lock.LockResource(ruleID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateAffinityRuleMemberRequest{
		AffinityRuleID: ruleID,
		InstanceID:     d.Get("instance_id").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateAffinityRuleMemberRequest: %+v", createReq))

	tflog.Info(ctx, "Creating AffinityRuleMember")
	taskRef, err := service.CreateAffinityRuleMember(createReq)
	if err != nil {
		return diag.Errorf("Error creating affinity rule member: %s", err)
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
		return diag.Errorf("Error waiting for affinity rule member with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceAffinityRuleMemberRead(ctx, d, meta)
}

func resourceAffinityRuleMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving affinity rule member", map[string]interface{}{
		"id": d.Id(),
	})
	arm, err := service.GetAffinityRuleMember(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleMemberNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("instance_id", arm.InstanceID)
	d.Set("affinity_rule_id", arm.AffinityRuleID)

	return nil
}

func resourceAffinityRuleMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceAffinityRuleMemberRead(ctx, d, meta)
}

func resourceAffinityRuleMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ruleID := d.Get("affinity_rule_id").(string)

	unlock := lock.LockResource(ruleID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing affinity rule member", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteAffinityRuleMember(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleMemberNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing affinity rule member with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for affinity rule member with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
