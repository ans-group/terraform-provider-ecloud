package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/ptr"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallPolicyCreate,
		Read:   resourceFirewallPolicyRead,
		Update: resourceFirewallPolicyUpdate,
		Delete: resourceFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			},
		},
	}
}

func resourceFirewallPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateFirewallPolicyRequest{
		RouterID: d.Get("router_id").(string),
		Sequence: d.Get("sequence").(int),
		Name:     d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateFirewallPolicyRequest: %+v", createReq)

	log.Print("[INFO] Creating firewall policy")
	policyID, err := service.CreateFirewallPolicy(createReq)
	if err != nil {
		return fmt.Errorf("Error creating firewall policy: %s", err)
	}

	d.SetId(policyID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, policyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", policyID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceFirewallPolicyRead(d, meta)
}

func resourceFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[DEBUG] Retrieving FirewallPolicy with ID [%s]", d.Id())
	policy, err := service.GetFirewallPolicy(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FirewallPolicyNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("router_id", policy.RouterID)
	d.Set("sequence", policy.Sequence)
	d.Set("name", policy.Name)

	return nil
}

func resourceFirewallPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[INFO] Updating firewall policy with ID [%s]", d.Id())
		err := service.PatchFirewallPolicy(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating firewall policy with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FirewallPolicySyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for firewall policy with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceFirewallPolicyRead(d, meta)
}

func resourceFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing firewall policy with ID [%s]", d.Id())
	err := service.DeleteFirewallPolicy(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing firewall policy with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    FirewallPolicySyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for firewall policy with ID [%s] to be deleted: %s", d.Id(), err)
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
