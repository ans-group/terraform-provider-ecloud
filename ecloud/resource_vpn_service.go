package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVPNService() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPNServiceCreate,
		Read:   resourceVPNServiceRead,
		Update: resourceVPNServiceUpdate,
		Delete: resourceVPNServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
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

func resourceVPNServiceCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPNServiceRequest{
		RouterID: d.Get("router_id").(string),
		Name:     d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateVPNServiceRequest: %+v", createReq)

	log.Print("[INFO] Creating VPN service")
	taskRef, err := service.CreateVPNService(createReq)
	if err != nil {
		return fmt.Errorf("Error creating VPN service: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPN service with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
	}

	return resourceVPNServiceRead(d, meta)
}

func resourceVPNServiceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving VPNService with ID [%s]", d.Id())
	vpc, err := service.GetVPNService(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNServiceNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("router_id", vpc.RouterID)
	d.Set("name", vpc.Name)

	return nil
}

func resourceVPNServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPNServiceRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating VPNService with ID [%s]", d.Id())
		taskRef, err := service.PatchVPNService(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating VPNService with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for VPN service with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceVPNServiceRead(d, meta)
}

func resourceVPNServiceDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing VPNService with ID [%s]", d.Id())
	taskID, err := service.DeleteVPNService(d.Id())
	if err != nil {
		return fmt.Errorf("Error VPNService with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPNService with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
