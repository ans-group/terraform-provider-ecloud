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

func resourceNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkRuleCreate,
		Read:   resourceNetworkRuleRead,
		Update: resourceNetworkRuleUpdate,
		Delete: resourceNetworkRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_policy_id": {
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
				Required: true,
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
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkRuleCreate(d *schema.ResourceData, meta interface{}) error {
	networkPolicyID := d.Get("network_policy_id").(string)
	unlock := lock.LockResource(networkPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	portsExpanded, err := expandCreateNetworkRuleRequestPorts(d.Get("port").([]interface{}))
	if err != nil {
		return err
	}

	createReq := ecloudservice.CreateNetworkRuleRequest{
		NetworkPolicyID: networkPolicyID,
		Name:            d.Get("name").(string),
		Sequence:        d.Get("sequence").(int),
		Source:          d.Get("source").(string),
		Destination:     d.Get("destination").(string),
		Enabled:         d.Get("enabled").(bool),
		Ports:           portsExpanded,
	}

	direction := d.Get("direction").(string)
	directionParsed, err := ecloudservice.ParseNetworkRuleDirection(direction)
	if err != nil {
		return err
	}
	createReq.Direction = directionParsed

	action := d.Get("action").(string)
	actionParsed, err := ecloudservice.ParseNetworkRuleAction(action)
	if err != nil {
		return err
	}
	createReq.Action = actionParsed

	log.Printf("[DEBUG] Created CreateNetworkRuleRequest: %+v", createReq)

	log.Print("[INFO] Creating network rule")
	task, err := service.CreateNetworkRule(createReq)
	if err != nil {
		return fmt.Errorf("Error creating network rule: %s", err)
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
		return fmt.Errorf("Error waiting for network policy with ID [%s] to return task status of [%s]: %s", networkPolicyID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceNetworkRuleRead(d, meta)
}

func resourceNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving network rule with ID [%s]", d.Id())
	rule, err := service.GetNetworkRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NetworkRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	log.Printf("[INFO] Retrieving network rule ports for network rule with ID [%s]", d.Id())
	// ports, err := service.GetNetworkRuleNetworkRulePorts(d.Id(), connection.APIRequestParameters{})
	
	//using filter parameter in request until dedicated API endpoint is 
	//added for service.GetNetworkRuleNetworkRulePorts().
	params := connection.APIRequestParameters{}
	params.WithFilter(*connection.NewAPIRequestFiltering("network_rule_id", connection.EQOperator, []string{d.Id()}))
	
	ports, err := service.GetNetworkRulePorts(params)
	if err != nil {
		return err
	}

	d.Set("network_policy_id", rule.NetworkPolicyID)
	d.Set("name", rule.Name)
	d.Set("sequence", rule.Sequence)
	d.Set("source", rule.Source)
	d.Set("destination", rule.Destination)
	d.Set("action", rule.Action)
	d.Set("direction", rule.Direction)
	d.Set("enabled", rule.Enabled)
	d.Set("port", flattenNetworkRulePorts(ports))

	return nil
}

func resourceNetworkRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	networkPolicyID := d.Get("network_policy_id").(string)
	unlock := lock.LockResource(networkPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)
	hasChange := false

	patchReq := ecloudservice.PatchNetworkRuleRequest{}

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
		actionParsed, err := ecloudservice.ParseNetworkRuleAction(action)
		if err != nil {
			return err
		}

		patchReq.Action = actionParsed
	}

	if d.HasChange("direction") {
		hasChange = true

		direction := d.Get("direction").(string)
		directionParsed, err := ecloudservice.ParseNetworkRuleDirection(direction)
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

		portsExpanded, err := expandUpdateNetworkRuleRequestPorts(d.Get("port").([]interface{}))
		if err != nil {
			return err
		}

		patchReq.Ports = portsExpanded
	}

	if hasChange {
		log.Printf("[INFO] Updating network rule with ID [%s]", d.Id())
		task, err := service.PatchNetworkRule(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating firewall rule with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for network policy with ID [%s] to return task status of [%s]: %s", networkPolicyID, ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceNetworkRuleRead(d, meta)
}

func resourceNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	networkPolicyID := d.Get("network_policy_id").(string)
	unlock := lock.LockResource(networkPolicyID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing network rule with ID [%s]", d.Id())
	taskID, err := service.DeleteNetworkRule(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing network rule with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for network policy with ID [%s] to return task status of [%s]: %s", networkPolicyID, ecloudservice.TaskStatusComplete, err)
	}

	return nil
}
