package ecloud

import (
	"fmt"
	"log"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/terraform-provider-ecloud/pkg/lock"
)

func resourceAffinityRuleMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAffinityRuleMemberCreate,
		Read:   resourceAffinityRuleMemberRead,
		Update: resourceAffinityRuleMemberUpdate,
		Delete: resourceAffinityRuleMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceAffinityRuleMemberCreate(d *schema.ResourceData, meta interface{}) error {
	ruleID := d.Get("affinity_rule_id").(string)
	if len(ruleID) < 1 {
		return fmt.Errorf("Invalid affinity rule ID: %s", ruleID)
	}

	unlock := lock.LockResource(ruleID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateAffinityRuleMemberRequest{
		AffinityRuleID: ruleID,
		InstanceID: d.Get("instance_id").(string),
	}

	log.Printf("[DEBUG] Created CreateAffinityRuleMemberRequest: %+v", createReq)

	log.Print("[INFO] Creating AffinityRuleMember")
	taskRef, err := service.CreateAffinityRuleMember(createReq)
	if err != nil {
		return fmt.Errorf("Error creating affinity rule member: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for affinity rule member with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceAffinityRuleMemberRead(d, meta)
}

func resourceAffinityRuleMemberRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving affinity rule member with ID [%s]", d.Id())
	arm, err := service.GetAffinityRuleMember(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleMemberNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("instance_id", arm.InstanceID)
	d.Set("affinity_rule_id", arm.AffinityRuleID)

	return nil
}

func resourceAffinityRuleMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAffinityRuleMemberRead(d, meta)
}

func resourceAffinityRuleMemberDelete(d *schema.ResourceData, meta interface{}) error {
	ruleID := d.Get("affinity_rule_id").(string)

	unlock := lock.LockResource(ruleID)
	defer unlock()

	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing affinity rule member with ID [%s]", d.Id())
	taskID, err := service.DeleteAffinityRuleMember(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleMemberNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing affinity rule member with ID [%s]: %s", d.Id(), err)
		}
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
		return fmt.Errorf("Error waiting for affinity rule member with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
