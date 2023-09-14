package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNATOverloadRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNATOverloadRuleCreate,
		ReadContext:   resourceNATOverloadRuleRead,
		UpdateContext: resourceNATOverloadRuleUpdate,
		DeleteContext: resourceNATOverloadRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNATOverloadRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	actionParsed, err := ecloud.ParseNATOverloadRuleAction(d.Get("action").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := ecloudservice.CreateNATOverloadRuleRequest{
		NetworkID:    d.Get("network_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		Subnet:       d.Get("subnet").(string),
		Action:       actionParsed,
		Name:         d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateNATOverloadRuleRequest: %+v", createReq))

	tflog.Info(ctx, "Creating NAT overload rule")
	taskRef, err := service.CreateNATOverloadRule(createReq)
	if err != nil {
		return diag.Errorf("Error creating NAT overload rule: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceNATOverloadRuleRead(ctx, d, meta)
}

func resourceNATOverloadRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving NAT overload rule", map[string]interface{}{
		"id": d.Id(),
	})
	rule, err := service.GetNATOverloadRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NATOverloadRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", rule.Name)
	d.Set("network_id", rule.NetworkID)
	d.Set("floating_ip_id", rule.FloatingIPID)
	d.Set("subnet", rule.Subnet)
	d.Set("action", rule.Action)

	return nil
}

func resourceNATOverloadRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	patchReq := ecloudservice.PatchNATOverloadRuleRequest{}
	changed := false
	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
		changed = true
	}
	if d.HasChange("subnet") {
		patchReq.Subnet = d.Get("subnet").(string)
		changed = true
	}

	if changed {
		tflog.Info(ctx, "Updating NAT overload rule", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.PatchNATOverloadRule(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating network with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceNATOverloadRuleRead(ctx, d, meta)
}

func resourceNATOverloadRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing NAT overload rule", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteNATOverloadRule(d.Id())
	if err != nil {
		return diag.Errorf("Error removing network with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
	}

	return nil
}
