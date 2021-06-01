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

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
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
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vcpu_cores": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"volume_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"volume_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					iops := []int{300, 600, 1200, 2500}
					intInSlice := func(slice []int, value int) bool {
						for _, s := range slice {
							if s == value {
								return true
							}
						}
						return false
					}

					if !intInSlice(iops, v) {
						errs = append(errs, fmt.Errorf("%q must be a valid IOPS value [300, 600, 1200, 2500], got: %d", key, v))
					}
					return
				},
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"requires_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateInstanceRequest{
		VPCID:              d.Get("vpc_id").(string),
		Name:               d.Get("name").(string),
		ImageID:            d.Get("image_id").(string),
		UserScript:         d.Get("user_script").(string),
		VCPUCores:          d.Get("vcpu_cores").(int),
		RAMCapacity:        d.Get("ram_capacity").(int),
		VolumeCapacity:     d.Get("volume_capacity").(int),
		VolumeIOPS:         d.Get("volume_iops").(int),
		Locked:             d.Get("locked").(bool),
		BackupEnabled:      d.Get("backup_enabled").(bool),
		NetworkID:          d.Get("network_id").(string),
		FloatingIPID:       d.Get("floating_ip_id").(string),
		RequiresFloatingIP: d.Get("requires_floating_ip").(bool),
	}
	log.Printf("[DEBUG] Created CreateInstanceRequest: %+v", createReq)

	log.Print("[INFO] Creating Instance")
	instanceID, err := service.CreateInstance(createReq)
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	d.SetId(instanceID)

	_, err = waitForResourceState(
		ecloudservice.SyncStatusComplete.String(),
		InstanceSyncStatusRefreshFunc(service, instanceID),
		d.Timeout(schema.TimeoutCreate),
	)
	if err != nil {
		return fmt.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", instanceID, ecloudservice.SyncStatusComplete, err)
	}

	rawIDs, ok := d.GetOk("data_volume_ids")
	if ok {
		for _, rawID := range rawIDs.(*schema.Set).List() {
			volumeID := rawID.(string)
			if len(volumeID) < 1 {
				continue
			}

			req := ecloudservice.AttachDetachInstanceVolumeRequest{
				VolumeID: volumeID,
			}
			log.Printf("[DEBUG] Created AttachDetachInstanceVolumeRequest: %+v", req)

			log.Printf("[INFO] Attaching volume to instance ID %s", d.Id())
			taskID, err := service.AttachInstanceVolume(d.Id(), req)
			if err != nil {
				return fmt.Errorf("Error attaching volume: %s", err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				InstanceTaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to be attached: %s", volumeID, err)
			}
		}
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving instance with ID [%s]", d.Id())
	instance, err := service.GetInstance(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.InstanceNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	volumes, err := service.GetInstanceVolumes(d.Id(), connection.APIRequestParameters{})
	if err != nil {
		return fmt.Errorf("Failed to retrieve instance volumes: %w", err)
	}

	//check we have 1 os volume
	var osVolume []ecloudservice.Volume
	for _, volume := range volumes {
		if volume.Type == ecloudservice.VolumeTypeOS {
			osVolume = append(osVolume, volume)
		}
	}
	if len(osVolume) != 1 {
		return fmt.Errorf("Unexpected number of OS volumes (%d), expected 1", len(volumes))
	}

	d.Set("vpc_id", instance.VPCID)
	d.Set("name", instance.Name)
	d.Set("image_id", instance.ImageID)
	d.Set("vcpu_cores", instance.VCPUCores)
	d.Set("ram_capacity", instance.RAMCapacity)
	d.Set("locked", instance.Locked)
	d.Set("backup_enabled", instance.BackupEnabled)
	d.Set("volume_capacity", osVolume[0].Capacity)
	d.Set("volume_iops", osVolume[0].IOPS)

	if d.Get("volume_id").(string) == "" {
		d.Set("volume_id", osVolume[0].ID)
	}

	d.Set("data_volume_ids", flattenInstanceDataVolumes(volumes))

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	patchReq := ecloudservice.PatchInstanceRequest{}
	hasChange := false
	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}
	if d.HasChange("vcpu_cores") {
		hasChange = true
		patchReq.VCPUCores = d.Get("vcpu_cores").(int)
	}
	if d.HasChange("ram_capacity") {
		hasChange = true
		patchReq.RAMCapacity = d.Get("ram_capacity").(int)
	}

	if hasChange {
		log.Printf("[INFO] Updating instance with ID [%s]", d.Id())
		err := service.PatchInstance(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating instance with ID [%s]: %w", d.Id(), err)
		}

		_, err = waitForResourceState(
			ecloudservice.SyncStatusComplete.String(),
			InstanceSyncStatusRefreshFunc(service, d.Id()),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	if d.HasChange("volume_capacity") {
		osVolumeID := d.Get("volume_id").(string)
		log.Printf("[INFO] Updating volume with ID [%s]", osVolumeID)
		taskID, err := service.PatchVolume(osVolumeID, ecloudservice.PatchVolumeRequest{
			Capacity: d.Get("volume_capacity").(int),
		})
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", osVolumeID, err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			VolumeTaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", osVolumeID, ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("volume_iops") {
		osVolumeID := d.Get("volume_id").(string)
		log.Printf("[INFO] Updating volume with ID [%s]", osVolumeID)
		taskID, err := service.PatchVolume(osVolumeID, ecloudservice.PatchVolumeRequest{
			IOPS: d.Get("volume_iops").(int),
		})
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", osVolumeID, err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			VolumeTaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", osVolumeID, ecloudservice.TaskStatusComplete, err)
		}
	}

	//manage attached data volumes
	if d.HasChange("data_volume_ids") {
		oldRaw, newRaw := d.GetChange("data_volume_ids")

		oldIDs := oldRaw.(*schema.Set).List()
		newIDs := newRaw.(*schema.Set).List()

		//diff new against old for attach
		for _, id := range newIDs {
			volumeID := id.(string)
			if rawVolumeExistsById(oldIDs, volumeID) {
				continue
			}

			attachReq := ecloudservice.AttachDetachInstanceVolumeRequest{
				VolumeID: volumeID,
			}
			log.Printf("[DEBUG] Created AttachDetachInstanceVolumeRequest: %+v", attachReq)

			log.Printf("[INFO] Attaching volume to instance ID [%s]", d.Id())
			taskID, err := service.AttachInstanceVolume(d.Id(), attachReq)
			if err != nil {
				return fmt.Errorf("Error attaching volume: %s", err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				InstanceTaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", volumeID, ecloudservice.TaskStatusComplete, err)
			}
		}

		//diff old against new for detach
		for _, id := range oldIDs {
			volumeID := id.(string)
			if rawVolumeExistsById(newIDs, volumeID) {
				continue
			}

			detachReq := ecloudservice.AttachDetachInstanceVolumeRequest{
				VolumeID: volumeID,
			}
			log.Printf("[DEBUG] Created DetachVolumeRequest: %+v", detachReq)

			log.Printf("[INFO] Detaching volume from instance ID [%s]", d.Id())
			taskID, err := service.DetachInstanceVolume(d.Id(), detachReq)
			if err != nil {
				return fmt.Errorf("Error detaching volume: %s", err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				InstanceTaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", volumeID, ecloudservice.TaskStatusComplete, err)
			}
		}
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing instance with ID [%s]", d.Id())
	err := service.DeleteInstance(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing instance with ID [%s]: %s", d.Id(), err)
	}

	_, err = waitForResourceState(
		"Deleted",
		InstanceSyncStatusRefreshFunc(service, d.Id()),
		d.Timeout(schema.TimeoutDelete),
	)
	if err != nil {
		return fmt.Errorf("Error waiting for instance with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// InstanceSyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func InstanceSyncStatusRefreshFunc(service ecloudservice.ECloudService, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := service.GetInstance(instanceID)
		if err != nil {
			if _, ok := err.(*ecloudservice.InstanceNotFoundError); ok {
				return instance, "Deleted", nil
			}
			return nil, "", err
		}

		if instance.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update instance - review logs")
		}

		return instance, instance.Sync.Status.String(), nil
	}
}

// InstanceTaskStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func InstanceTaskStatusRefreshFunc(service ecloudservice.ECloudService, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Retrieving task status for taskID: [%s]", taskID)

		task, err := service.GetTask(taskID)
		if err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] TaskID: %s has status: %s", task.ID, task.Status)

		if task.Status == ecloudservice.TaskStatusFailed {
			return nil, "", fmt.Errorf("Task with ID: %s has status of %s", task.ID, task.Status)
		}

		return "", task.Status.String(), nil
	}
}

// waitForResourceState is a wrapper for the resource.StateChangeConf helper in order to reduce duplication
func waitForResourceState(targetState string, refreshFunc resource.StateRefreshFunc, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Target:     []string{targetState},
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}
