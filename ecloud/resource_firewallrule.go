package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
	"github.com/ukfast/terraform-provider-ecloud/pkg/lock"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"firewall_policy_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	firewallPolicyID := d.Get("firewall_policy_id").(string)
	unlock := lock.LockResource(firewallPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	portsExpanded, err := expandCreateFirewallRuleRequestPorts(d.Get("port").([]interface{}))
	if err != nil {
		return err
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
	directionParsed, err := ecloudservice.ParseFirewallRuleDirection(direction)
	if err != nil {
		return err
	}
	createReq.Direction = directionParsed

	action := d.Get("action").(string)
	actionParsed, err := ecloudservice.ParseFirewallRuleAction(action)
	if err != nil {
		return err
	}
	createReq.Action = actionParsed

	log.Printf("[DEBUG] Created CreateFirewallRuleRequest: %+v", createReq)

	log.Print("[INFO] Creating firewall rule")
	ruleID, err := service.CreateFirewallRule(createReq)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule: %s", err)
	}

	d.SetId(ruleID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, firewallPolicyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", firewallPolicyID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving firewall rule with ID [%s]", d.Id())
	rule, err := service.GetFirewallRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FirewallRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	log.Printf("[INFO] Retrieving firewall rule ports for firewall rule with ID [%s]", d.Id())
	ports, err := service.GetFirewallRuleFirewallRulePorts(d.Id(), connection.APIRequestParameters{})
	if err != nil {
		return err
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

func resourceFirewallRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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
		actionParsed, err := ecloudservice.ParseFirewallRuleAction(action)
		if err != nil {
			return err
		}

		patchReq.Action = actionParsed
	}

	if d.HasChange("action") {
		hasChange = true

		direction := d.Get("direction").(string)
		directionParsed, err := ecloudservice.ParseFirewallRuleDirection(direction)
		if err != nil {
			return err
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
			return err
		}

		patchReq.Ports = portsExpanded
	}

	if hasChange {
		log.Printf("[INFO] Updating firewall rule with ID [%s]", d.Id())
		err := service.PatchFirewallRule(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating firewall rule with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FirewallPolicySyncStatusRefreshFunc(service, firewallPolicyID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", firewallPolicyID, ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	firewallPolicyID := d.Get("firewall_policy_id").(string)
	unlock := lock.LockResource(firewallPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing firewall rule with ID [%s]", d.Id())
	err := service.DeleteFirewallRule(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing firewall rule with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target: []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			rule, err := service.GetFirewallRule(d.Id())
			if err != nil {
				if _, ok := err.(*ecloudservice.FirewallRuleNotFoundError); ok {
					return rule, "Deleted", nil
				}
				return nil, "", err
			}

			return rule, "", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for firewall rule with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
