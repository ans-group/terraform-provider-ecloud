package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallPolicyCreate,
		ReadContext:   resourceFirewallPolicyRead,
		UpdateContext: resourceFirewallPolicyUpdate,
		DeleteContext: resourceFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceFirewallPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateFirewallPolicyRequest{
		RouterID: d.Get("router_id").(string),
		Sequence: d.Get("sequence").(int),
		Name:     d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateFirewallPolicyRequest: %+v", createReq))

	tflog.Info(ctx, "Creating firewall policy")
	policy, err := service.CreateFirewallPolicy(createReq)
	if err != nil {
		return diag.Errorf("Error creating firewall policy: %s", err)
	}

	d.SetId(policy.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, policy.ResourceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", policy.ResourceID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceFirewallPolicyRead(ctx, d, meta)
}

func resourceFirewallPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving firewall policy", map[string]interface{}{
		"id": d.Id(),
	})
	policy, err := service.GetFirewallPolicy(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FirewallPolicyNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("router_id", policy.RouterID)
	d.Set("sequence", policy.Sequence)
	d.Set("name", policy.Name)

	return nil
}

func resourceFirewallPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)
	hasChange := false

	patchReq := ecloudservice.PatchFirewallPolicyRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("sequence") {
		hasChange = true
		patchReq.Sequence = ptr.Int(d.Get("sequence").(int))
	}

	if hasChange {
		tflog.Info(ctx, "Updating firewall policy", map[string]interface{}{
			"id": d.Id(),
		})
		_, err := service.PatchFirewallPolicy(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating firewall policy with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FirewallPolicySyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceFirewallPolicyRead(ctx, d, meta)
}

func resourceFirewallPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing firewall policy", map[string]interface{}{
		"id": d.Id(),
	})
	_, err := service.DeleteFirewallPolicy(d.Id())
	if err != nil {
		return diag.Errorf("Error removing firewall policy with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for firewall policy with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// FirewallPolicySyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func FirewallPolicySyncStatusRefreshFunc(service ecloudservice.ECloudService, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := service.GetFirewallPolicy(policyID)
		if err != nil {
			if _, ok := err.(*ecloudservice.FirewallPolicyNotFoundError); ok {
				return policy, "Deleted", nil
			}
			return nil, "", err
		}

		if policy.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update firewall policy - review logs")
		}

		return policy, policy.Sync.Status.String(), nil
	}
}
