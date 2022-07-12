package ecloud

import (
	"fmt"
	"log"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostGroupCreate,
		Read:   resourceHostGroupRead,
		Update: resourceHostGroupUpdate,
		Delete: resourceHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_spec_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"windows_enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"availability_zone_id": {
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

func resourceHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateHostGroupRequest{
		VPCID:              d.Get("vpc_id").(string),
		HostSpecID:         d.Get("host_spec_id").(string),
		WindowsEnabled:     d.Get("windows_enabled").(bool),
		Name:               d.Get("name").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}
	log.Printf("[DEBUG] Created CreateHostGroupRequest: %+v", createReq)

	log.Print("[INFO] Creating Host Group")
	task, err := service.CreateHostGroup(createReq)
	if err != nil {
		return fmt.Errorf("Error creating host group: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for host group with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceHostGroupRead(d, meta)
}

func resourceHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving host group with ID [%s]", d.Id())
	hostGroup, err := service.GetHostGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostGroupNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", hostGroup.VPCID)
	d.Set("host_spec_id", hostGroup.HostSpecID)
	d.Set("availability_zone_id", hostGroup.AvailabilityZoneID)
	d.Set("windows_enabled", hostGroup.WindowsEnabled)
	d.Set("name", hostGroup.Name)

	return nil
}

func resourceHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating host group with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchHostGroupRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchHostGroup(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating host group with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for host group with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceHostGroupRead(d, meta)
}

func resourceHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing host group with ID [%s]", d.Id())
	taskID, err := service.DeleteHostGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostGroupNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing host group with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for host group with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
