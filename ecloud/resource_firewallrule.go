package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/terraform-provider-ecloud/pkg/lock"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallRuleCreate,
		ReadContext:   resourceFirewallRuleRead,
		UpdateContext: resourceFirewallRuleUpdate,
		DeleteContext: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"firewall_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	firewallPolicyID := d.Get("firewall_policy_id").(string)
	unlock := lock.LockResource(firewallPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	portsExpanded, err := expandCreateFirewallRuleRequestPorts(d.Get("port").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := ecloudservice.CreateFirewallRuleRequest{
		FirewallPolicyID: firewallPolicyID,
		Name:             d.Get("name").(string),
		Sequence:         d.Get("sequence").(int),
		Source:           d.Get("source").(string),
		Destination:      d.Get("destination").(string),
		Enabled:          d.Get("enabled").(bool),
		Ports:            portsExpanded,
	}

	direction := d.Get("direction").(string)
	directionParsed, err := ecloudservice.FirewallRuleDirectionEnum.Parse(direction)
	if err != nil {
		return diag.FromErr(err)
	}
	createReq.Direction = directionParsed

	action := d.Get("action").(string)
	actionParsed, err := ecloudservice.FirewallRuleActionEnum.Parse(action)
	if err != nil {
		return diag.FromErr(err)
	}
	createReq.Action = actionParsed

	tflog.Debug(ctx, fmt.Sprintf("Created CreateFirewallRuleRequest: %+v", createReq))

	tflog.Info(ctx, "Creating firewall rule")
	rule, err := service.CreateFirewallRule(createReq)
	if err != nil {
		return diag.Errorf("Error creating firewall rule: %s", err)
	}

	d.SetId(rule.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, firewallPolicyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", firewallPolicyID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceFirewallRuleRead(ctx, d, meta)
}

func resourceFirewallRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving firewall rule", map[string]interface{}{
		"id": d.Id(),
	})
	rule, err := service.GetFirewallRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FirewallRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	tflog.Info(ctx, "Retrieving firewall rule ports", map[string]interface{}{
		"id": d.Id(),
	})
	ports, err := service.GetFirewallRuleFirewallRulePorts(d.Id(), connection.APIRequestParameters{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("firewall_policy_id", rule.FirewallPolicyID)
	d.Set("name", rule.Name)
	d.Set("sequence", rule.Sequence)
	d.Set("source", rule.Source)
	d.Set("destination", rule.Destination)
	d.Set("action", rule.Action)
	d.Set("direction", rule.Direction)
	d.Set("enabled", rule.Enabled)
	d.Set("port", flattenFirewallRulePorts(ports))

	return nil
}

func resourceFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	firewallPolicyID := d.Get("firewall_policy_id").(string)
	unlock := lock.LockResource(firewallPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)
	hasChange := false

	patchReq := ecloudservice.PatchFirewallRuleRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("sequence") {
		hasChange = true
		patchReq.Sequence = ptr.Int(d.Get("sequence").(int))
	}

	if d.HasChange("source") {
		hasChange = true
		patchReq.Source = d.Get("source").(string)
	}

	if d.HasChange("destination") {
		hasChange = true
		patchReq.Destination = d.Get("destination").(string)
	}

	if d.HasChange("action") {
		hasChange = true

		action := d.Get("action").(string)
		actionParsed, err := ecloudservice.FirewallRuleActionEnum.Parse(action)
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Action = actionParsed
	}

	if d.HasChange("direction") {
		hasChange = true

		direction := d.Get("direction").(string)
		directionParsed, err := ecloudservice.FirewallRuleDirectionEnum.Parse(direction)
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Direction = directionParsed
	}

	if d.HasChange("enabled") {
		hasChange = true
		patchReq.Enabled = ptr.Bool(d.Get("enabled").(bool))
	}

	if d.HasChange("port") {
		hasChange = true

		portsExpanded, err := expandUpdateFirewallRuleRequestPorts(d.Get("port").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		patchReq.Ports = portsExpanded
	}

	if hasChange {
		tflog.Info(ctx, "Updating firewall rule", map[string]interface{}{
			"id": d.Id(),
		})

		_, err := service.PatchFirewallRule(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating firewall rule with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FirewallPolicySyncStatusRefreshFunc(service, firewallPolicyID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", firewallPolicyID, ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceFirewallRuleRead(ctx, d, meta)
}

func resourceFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	firewallPolicyID := d.Get("firewall_policy_id").(string)
	unlock := lock.LockResource(firewallPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing firewall rule", map[string]interface{}{
		"id": d.Id(),
	})
	_, err := service.DeleteFirewallRule(d.Id())
	if err != nil {
		return diag.Errorf("Error removing firewall rule with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, firewallPolicyID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", firewallPolicyID, ecloudservice.SyncStatusComplete, err)
	}

	return nil
}
