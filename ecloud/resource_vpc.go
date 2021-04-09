package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCCreate,
		Read:   resourceVPCRead,
		Update: resourceVPCUpdate,
		Delete: resourceVPCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVPCCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPCRequest{
		RegionID: d.Get("region_id").(string),
		Name:     d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateVPCRequest: %+v", createReq)

	log.Print("[INFO] Creating VPC")
	vpcID, err := service.CreateVPC(createReq)
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpcID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    VPCSyncStatusRefreshFunc(service, vpcID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPC with ID [%s] to return sync status of [%s]: %s", vpcID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceVPCRead(d, meta)
}

func resourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving VPC with ID [%s]", d.Id())
	vpc, err := service.GetVPC(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPCNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("region_id", vpc.RegionID)
	d.Set("name", vpc.Name)

	return nil
}

func resourceVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPCRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating VPC with ID [%s]", d.Id())
		err := service.PatchVPC(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating VPC with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    VPCSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for VPC with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceVPCRead(d, meta)
}

func resourceVPCDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing VPC with ID [%s]", d.Id())
	err := service.DeleteVPC(d.Id())
	if err != nil {
		return fmt.Errorf("Error VPC with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    VPCSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPC with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// VPCSyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func VPCSyncStatusRefreshFunc(service ecloudservice.ECloudService, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpc, err := service.GetVPC(vpcID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPCNotFoundError); ok {
				return vpc, "Deleted", nil
			}
			return nil, "", err
		}

		if vpc.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update VPC - review logs")
		}

		return vpc, vpc.Sync.Status.String(), nil
	}
}
