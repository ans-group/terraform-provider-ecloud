package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceAffinityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAffinityRuleCreate,
		Read:   resourceAffinityRuleRead,
		Update: resourceAffinityRuleUpdate,
		Delete: resourceAffinityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := ecloudservice.ParseAffinityRuleType(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid affinity rule type [affinity, anti-affinity], got: %s", key, v))
					}
					return
				},
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceAffinityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	ruleType, err := ecloudservice.ParseAffinityRuleType(d.Get("type").(string))
	if err != nil {
		return fmt.Errorf("Error parsing affinity rule type: %s", err)
	}

	createReq := ecloudservice.CreateAffinityRuleRequest{
		VPCID:              d.Get("vpc_id").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
		Name:               d.Get("name").(string),
		Type:               ruleType,
	}

	log.Printf("[DEBUG] Created CreateAffinityRuleRequest: %+v", createReq)

	log.Print("[INFO] Creating AffinityRule")
	taskRef, err := service.CreateAffinityRule(createReq)
	if err != nil {
		return fmt.Errorf("Error creating affinity rule: %s", err)
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
		return fmt.Errorf("Error waiting for affinity rule with ID [%s] to be created: %s", d.Id(), err)
	}

	//add members to rule
	if rawIDs, ok := d.GetOk("instance_ids"); ok {
		for _, rawID := range rawIDs.(*schema.Set).List() {
			memberID := rawID.(string)
			if len(memberID) < 1 {
				continue
			}

			req := ecloudservice.CreateAffinityRuleMemberRequest{
				InstanceID: memberID,
			}

			log.Printf("[DEBUG] Created CreateAffinityRuleMemberRequest: %+v", req)

			log.Printf("[INFO] Adding instance member [%s] to affinity rule ID [%s]", memberID, d.Id())
			taskRef, err := service.CreateAffinityRuleMember(d.Id(), req)
			if err != nil {
				return fmt.Errorf("Error creating affinity rule member: %s", err)
			}

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
		}
	}

	return resourceAffinityRuleRead(d, meta)
}

func resourceAffinityRuleRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving affinity rule with ID [%s]", d.Id())
	ar, err := service.GetAffinityRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	ruleMembers, err := service.GetAffinityRuleMembers(d.Id(), connection.APIRequestParameters{})
	if err != nil {
		return fmt.Errorf("Error retrieving affinity rule members for rule ID [%s]: %w", d.Id(), err)
	}

	d.Set("vpc_id", ar.VPCID)
	d.Set("name", ar.Name)
	d.Set("availability_zone_id", ar.AvailabilityZoneID)
	d.Set("type", ar.Type)
	d.Set("instance_ids", flattenAffinityRuleMembers(ruleMembers))

	return nil
}

func resourceAffinityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating affinity rule with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchAffinityRuleRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchAffinityRule(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating affinity rule with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for affinity rule with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("instance_ids") {
		oldRaw, newRaw := d.GetChange("instance_ids")

		oldIDs := oldRaw.(*schema.Set).List()
		newIDs := newRaw.(*schema.Set).List()

		ruleMembers, err := service.GetAffinityRuleMembers(d.Id(), connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Error retrieving affinity rule members for rule ID [%s]: %w", d.Id(), err)
		}
		getMemberByInstanceID := func(slice []ecloudservice.AffinityRuleMember, value string) string {
			for _, s := range slice {
				if s.InstanceID == value {
					return s.ID
				}
			}
			return ""
		}

		for _, id := range oldIDs {
			instanceID := id.(string)
			if rawMemberExistsById(newIDs, instanceID) {
				continue
			}

			ruleMemberID := getMemberByInstanceID(ruleMembers, instanceID)

			log.Printf("[INFO] Removing instance member [%s] from affinity rule ID [%s]", instanceID, d.Id())

			taskID, err := service.DeleteAffinityRuleMember(d.Id(), ruleMemberID)
			if err != nil {
				return fmt.Errorf("Error deleting affinity rule member: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(service, taskID),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for affinity rule member for instance ID [%s] to be delete: %s", instanceID, err)
			}
		}

		for _, id := range newIDs {
			instanceID := id.(string)
			if rawMemberExistsById(oldIDs, instanceID) {
				continue
			}

			req := ecloudservice.CreateAffinityRuleMemberRequest{
				InstanceID: instanceID,
			}

			log.Printf("[DEBUG] Created CreateAffinityRuleMemberRequest: %+v", req)

			log.Printf("[INFO] Adding instance member [%s] to affinity rule ID [%s]", instanceID, d.Id())
			taskRef, err := service.CreateAffinityRuleMember(d.Id(), req)
			if err != nil {
				return fmt.Errorf("Error creating affinity rule member: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      10 * time.Second,
				MinTimeout: 5 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for affinity rule member for instance ID [%s] to be created: %s", instanceID, err)
			}

		}
	}

	return resourceAffinityRuleRead(d, meta)
}

func resourceAffinityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	//delete members first if they exist
	if memberInstances, ok := d.GetOk("instance_ids"); ok {
		for _, memberInstance := range memberInstances.(*schema.Set).List() {
			instanceID := memberInstance.(string)

			if len(instanceID) < 1 {
				continue
			}

			// get affinity rule member ID by instanceID
			params := connection.APIRequestParameters{}
			params.WithFilter(*connection.NewAPIRequestFiltering("instance_id", connection.EQOperator, []string{instanceID}))

			arMembers, err := service.GetAffinityRuleMembers(d.Id(), params)
			if err != nil {
				return fmt.Errorf("Error retrieving affinity rule members for rule ID [%s]: %s", d.Id(), err)
			}

			if len(arMembers) < 1 {
				//resource may have been deleted already, so skip
				continue
			}

			if len(arMembers) > 1 {
				return fmt.Errorf("More than 1 affinity rule member found for instance ID [%s]", instanceID)
			}

			log.Printf("[INFO] Removing affinity rule member ID [%s] for instance ID [%s]", arMembers[0].ID, instanceID)
			taskID, err := service.DeleteAffinityRuleMember(d.Id(), arMembers[0].ID)
			if err != nil {
				switch err.(type) {
				case *ecloudservice.AffinityRuleMemberNotFoundError:
					log.Printf("[DEBUG] Affinity rule member ID [%s] for instance ID [%s] not found. Skipping delete.", arMembers[0].ID, instanceID)
				default:
					return fmt.Errorf("Error removing affinity rule member ID [%s] for instance ID [%s]: %s", arMembers[0].ID, instanceID, err)
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
				return fmt.Errorf("Error waiting for affinity rule member ID [%s] for instance ID [%s] to be removed: %w", arMembers[0].ID, instanceID, err)
			}
		}
	}

	log.Printf("[INFO] Removing affinity rule with ID [%s]", d.Id())
	taskID, err := service.DeleteAffinityRule(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.AffinityRuleNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing affinity rule with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for affinity rule with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}