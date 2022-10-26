package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNATOverloadRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNATOverloadRuleCreate,
		Read:   resourceNATOverloadRuleRead,
		Update: resourceNATOverloadRuleUpdate,
		Delete: resourceNATOverloadRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceNATOverloadRuleCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	actionParsed, err := ecloud.ParseNATOverloadRuleAction(d.Get("action").(string))
	if err != nil {
		return err
	}

	createReq := ecloudservice.CreateNATOverloadRuleRequest{
		NetworkID:    d.Get("network_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		Subnet:       d.Get("subnet").(string),
		Action:       actionParsed,
		Name:         d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateNATOverloadRuleRequest: %+v", createReq)

	log.Print("[INFO] Creating NAT overload rule")
	taskRef, err := service.CreateNATOverloadRule(createReq)
	if err != nil {
		return fmt.Errorf("Error creating NAT overload rule: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceNATOverloadRuleRead(d, meta)
}

func resourceNATOverloadRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving NAT overload rule with ID [%s]", d.Id())
	rule, err := service.GetNATOverloadRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NATOverloadRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", rule.Name)
	d.Set("network_id", rule.NetworkID)
	d.Set("floating_ip_id", rule.FloatingIPID)
	d.Set("subnet", rule.Subnet)
	d.Set("action", rule.Action)

	return nil
}

func resourceNATOverloadRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[INFO] Updating NAT overload rule with ID [%s]", d.Id())
		taskRef, err := service.PatchNATOverloadRule(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating network with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceNATOverloadRuleRead(d, meta)
}

func resourceNATOverloadRuleDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing network with ID [%s]", d.Id())
	taskID, err := service.DeleteNATOverloadRule(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing network with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
	}

	return nil
}
