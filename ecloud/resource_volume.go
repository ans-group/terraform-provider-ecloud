package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,
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
			},
			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVolumeRequest{
		VPCID:    d.Get("vpc_id").(string),
		Name:     d.Get("name").(string),
		Capacity: d.Get("capacity").(int),
		IOPS:     d.Get("iops").(int),
	}

	log.Printf("[DEBUG] Created CreateVolumeRequest: %+v", createReq)

	log.Print("[INFO] Creating Volume")
	taskRef, err := service.CreateVolume(createReq)
	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    VolumeSyncStatusRefreshFunc(service, d.Id(), taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for volume with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving volume with ID [%s]", d.Id())
	volume, err := service.GetVolume(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", volume.VPCID)
	d.Set("name", volume.Name)
	d.Set("capacity", volume.Capacity)
	d.Set("iops", volume.IOPS)

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	patchReq := ecloudservice.PatchVolumeRequest{}
	hasChange := false
	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}
	if d.HasChange("capacity") {
		hasChange = true
		patchReq.Capacity = d.Get("capacity").(int)
	}
	if d.HasChange("iops") {
		hasChange = true
		patchReq.IOPS = d.Get("iops").(int)
	}

	if hasChange {
		log.Printf("[INFO] Updating volume with ID [%s]", d.Id())
		taskID, err := service.PatchVolume(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    VolumeSyncStatusRefreshFunc(service, d.Id(), taskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}
	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing volume with ID [%s]", d.Id())
	taskID, err := service.DeleteVolume(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing volume with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    VolumeSyncStatusRefreshFunc(service, d.Id(), taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for volume with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

func VolumeSyncStatusRefreshFunc(service ecloudservice.ECloudService, volumeID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		//check resource exists first
		volume, err := service.GetVolume(volumeID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VolumeNotFoundError); ok {
				return volume, "Deleted", nil
			}
			return nil, "", err
		}

		//check task status on resource
		log.Printf("[DEBUG] Retrieving task status for taskID: [%s]", taskID)
		tasks, err := service.GetVolumeTasks(volumeID, *connection.NewAPIRequestParameters().WithFilter(
			connection.APIRequestFiltering{
				Property: "id",
				Operator: connection.EQOperator,
				Value: []string{taskID},
			},
		))
		if err != nil {
			return nil, "", err
		}
		if len(tasks) != 1 {
			return nil, "", fmt.Errorf("Expected 1 task, got %d", len(tasks))
		}

		log.Printf("[DEBUG] TaskID: %s has status: %s", tasks[0].ID, tasks[0].Status)

		if tasks[0].Status == ecloudservice.TaskStatusFailed {
			return nil, "", fmt.Errorf("Task with ID: %s has status of %s", tasks[0].ID, tasks[0].Status)
		}

		return volume, tasks[0].Status.String(), nil
	}
}
