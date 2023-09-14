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

func resourceNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPolicyCreate,
		ReadContext:   resourceNetworkPolicyRead,
		UpdateContext: resourceNetworkPolicyUpdate,
		DeleteContext: resourceNetworkPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catchall_rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ecloudservice.NetworkPolicyCatchallRuleActionReject.String(),
			},
			"catchall_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateNetworkPolicyRequest{
		NetworkID: d.Get("network_id").(string),
		Name:      d.Get("name").(string),
	}

	if catchallRuleAction, ok := d.GetOk("catchall_rule_action"); ok {
		action, err := ecloudservice.ParseNetworkPolicyCatchallRuleAction(catchallRuleAction.(string))
		if err != nil {
			return diag.Errorf("Error parsing network policy catch-all rule action: %s", err)
		}
		createReq.CatchallRuleAction = action
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateNetworkPolicyRequest: %+v", createReq))

	tflog.Info(ctx, "Creating network policy")
	task, err := service.CreateNetworkPolicy(createReq)
	if err != nil {
		return diag.Errorf("Error creating network policy: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for network policy with ID [%s] to return status of [%s]: %s", task.ResourceID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceNetworkPolicyRead(ctx, d, meta)
}

func resourceNetworkPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving Network Policy", map[string]interface{}{
		"id": d.Id(),
	})
	policy, err := service.GetNetworkPolicy(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NetworkPolicyNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("network_id", policy.NetworkID)
	d.Set("vpc_id", policy.VPCID)
	d.Set("name", policy.Name)

	if d.Get("catchall_rule_id").(string) == "" {
		tflog.Info(ctx, "Retrieving catchall rule for Network Policy", map[string]interface{}{
			"id": d.Id(),
		})

		params := connection.APIRequestParameters{}
		params.WithFilter(*connection.NewAPIRequestFiltering("type", connection.EQOperator, []string{"catchall"}))

		rules, err := service.GetNetworkPolicyNetworkRules(d.Id(), params)
		if err != nil {
			return diag.Errorf("Error retrieving network policy catch-all rule: %s", err)
		}

		if len(rules) > 1 {
			return diag.Errorf("More than 1 network policy catchall rule exists")
		}

		if len(rules) < 1 {
			return diag.Errorf("No catchall rule found for network policy")
		}

		d.Set("catchall_rule_id", rules[0].ID)
	}

	return nil
}

func resourceNetworkPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)
	hasChange := false

	patchReq := ecloudservice.PatchNetworkPolicyRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if hasChange {
		tflog.Info(ctx, "Updating Network Policy", map[string]interface{}{
			"id": d.Id(),
		})
		task, err := service.PatchNetworkPolicy(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating networking policy with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for network policy with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("catchall_rule_action") {
		// parse new rule action
		catchallRuleAction, err := ecloudservice.ParseNetworkPolicyCatchallRuleAction(d.Get("catchall_rule_action").(string))
		if err != nil {
			return diag.Errorf("Error parsing network rule action: %s", err)
		}

		// retrieve catchall rule by id
		rule, err := service.GetNetworkRule(d.Get("catchall_rule_id").(string))
		if err != nil {
			return diag.Errorf("Error retrieving network rule with ID [%s]: %s", d.Get("catchall_rule_id").(string), err)
		}

		// patch rule
		patchRuleReq := ecloudservice.PatchNetworkRuleRequest{
			Action: ecloudservice.NetworkRuleAction(catchallRuleAction),
		}
		tflog.Debug(ctx, fmt.Sprintf("Created PatchNetworkRuleRequest: %+v", patchRuleReq))

		task, err := service.PatchNetworkRule(rule.ID, patchRuleReq)
		if err != nil {
			return diag.Errorf("Error updating network rule action: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for network rule with ID [%s] to be deleted: %s", rule.ID, err)
		}
	}

	return resourceNetworkPolicyRead(ctx, d, meta)
}

func resourceNetworkPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing Network Policy", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteNetworkPolicy(d.Id())
	if err != nil {
		return diag.Errorf("Error removing network policy with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for network policy with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
