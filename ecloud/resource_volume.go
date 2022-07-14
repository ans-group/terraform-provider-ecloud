package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"availability_zone_id": {
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
			"volume_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVolumeRequest{
		VPCID:              d.Get("vpc_id").(string),
		Name:               d.Get("name").(string),
		Capacity:           d.Get("capacity").(int),
		IOPS:               d.Get("iops").(int),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}

	if volumeGroupID, ok := d.GetOk("volume_group_id"); ok {
		createReq.VolumeGroupID = volumeGroupID.(string)
	}

	log.Printf("[DEBUG] Created CreateVolumeRequest: %+v", createReq)

	log.Print("[INFO] Creating Volume")
	taskRef, err := service.CreateVolume(createReq)
	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for volume with ID [%s] to be created: %s", d.Id(), err)
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
	d.Set("availability_zone_id", volume.AvailabilityZoneID)
	d.Set("volume_group_id", volume.VolumeGroupID)

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
	if d.HasChange("volume_group_id") {
		hasChange = true
		patchReq.VolumeGroupID = ptr.String(d.Get("volume_group_id").(string))
	}

	if hasChange {
		log.Printf("[INFO] Updating volume with ID [%s]", d.Id())
		task, err := service.PatchVolume(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
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
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
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
