package ecloud

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkPolicyCreate,
		Read:   resourceNetworkPolicyRead,
		Update: resourceNetworkPolicyUpdate,
		Delete: resourceNetworkPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceNetworkPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateNetworkPolicyRequest{
		NetworkID: d.Get("network_id").(string),
		Name:      d.Get("name").(string),
	}

	if catchallRuleAction, ok := d.GetOk("catchall_rule_action"); ok {
		action, err := ecloudservice.ParseNetworkPolicyCatchallRuleAction(catchallRuleAction.(string))
		if err != nil {
			return fmt.Errorf("Error parsing network policy catch-all rule action: %s", err)
		}
		createReq.CatchallRuleAction = action
	}

	log.Printf("[DEBUG] Created CreateNetworkPolicyRequest: %+v", createReq)

	log.Print("[INFO] Creating network policy")
	task, err := service.CreateNetworkPolicy(createReq)
	if err != nil {
		return fmt.Errorf("Error creating network policy: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for network policy with ID [%s] to return status of [%s]: %s", task.ResourceID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceNetworkPolicyRead(d, meta)
}

func resourceNetworkPolicyRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[DEBUG] Retrieving Network Policy with ID [%s]", d.Id())
	policy, err := service.GetNetworkPolicy(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NetworkPolicyNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("network_id", policy.NetworkID)
	d.Set("vpc_id", policy.VPCID)
	d.Set("name", policy.Name)

	if d.Get("catchall_rule_id").(string) == "" {
		log.Printf("[DEBUG] Retrieving catchall rule for Network Policy with ID [%s]", d.Id())

		params := connection.APIRequestParameters{}
		params.WithFilter(*connection.NewAPIRequestFiltering("type", connection.EQOperator, []string{"catchall"}))

		rules, err := service.GetNetworkPolicyNetworkRules(d.Id(), params)
		if err != nil {
			return fmt.Errorf("Error retrieving network policy catch-all rule: %s", err)
		}

		if len(rules) > 1 {
			return errors.New("More than 1 network policy catchall rule exists")
		}

		if len(rules) < 1 {
			return errors.New("No catchall rule found for network policy")
		}

		d.Set("catchall_rule_id", rules[0].ID)
	}

	return nil
}

func resourceNetworkPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)
	hasChange := false

	patchReq := ecloudservice.PatchNetworkPolicyRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[INFO] Updating network policy with ID [%s]", d.Id())
		task, err := service.PatchNetworkPolicy(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating networking policy with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for network policy with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("catchall_rule_action") {
		//parse new rule action
		catchallRuleAction, err := ecloudservice.ParseNetworkPolicyCatchallRuleAction(d.Get("catchall_rule_action").(string))
		if err != nil {
			return fmt.Errorf("Error parsing network rule action: %s", err)
		}

		//retrieve catchall rule by id
		rule, err := service.GetNetworkRule(d.Get("catchall_rule_id").(string))
		if err != nil {
			return fmt.Errorf("Error retrieving network rule with ID [%s]: %s", d.Get("catchall_rule_id").(string), err)
		}

		//patch rule
		patchRuleReq := ecloudservice.PatchNetworkRuleRequest{
			Action: ecloudservice.NetworkRuleAction(catchallRuleAction),
		}
		log.Printf("[DEBUG] Created PatchNetworkRuleRequest: %+v", patchRuleReq)

		task, err := service.PatchNetworkRule(rule.ID, patchRuleReq)
		if err != nil {
			return fmt.Errorf("Error updating network rule action: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for network rule with ID [%s] to be deleted: %s", rule.ID, err)
		}
	}

	return resourceNetworkPolicyRead(d, meta)
}

func resourceNetworkPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing network policy with ID [%s]", d.Id())
	taskID, err := service.DeleteNetworkPolicy(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing network policy with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for network policy with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
