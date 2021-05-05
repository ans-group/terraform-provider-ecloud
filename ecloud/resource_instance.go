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
			"os_volume_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"os_volume_id": {
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
			},
			"requires_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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
		VolumeCapacity:     d.Get("os_volume_capacity").(int),
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

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    InstanceSyncStatusRefreshFunc(service, instanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
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

			attachVolumeRequest := ecloudservice.AttachVolumeRequest{
				InstanceID: instanceID,
			}
			log.Printf("[DEBUG] Created AttachVolumeRequest: %+v", attachVolumeRequest)

			log.Print("[INFO] Attaching volume")
			err = service.AttachVolume(volumeID, attachVolumeRequest)
			if err != nil {
				return fmt.Errorf("Error attaching volume: %s", err)
			}

			volStateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.SyncStatusComplete.String()},
				Refresh:    VolumeSyncStatusRefreshFunc(service, volumeID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = volStateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to be deleted: %s", volumeID, err)
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
	d.Set("os_volume_capacity", osVolume[0].Capacity)

	if d.Get("os_volume_id").(string) == "" {
		d.Set("os_volume_id", osVolume[0].ID)
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

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    InstanceSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	if d.HasChange("os_volume_capacity") {
		osVolumeID := d.Get("os_volume_id").(string)
		log.Printf("[INFO] Updating volume with ID [%s]", osVolumeID)
		service.PatchVolume(osVolumeID, ecloudservice.PatchVolumeRequest{
			Capacity: d.Get("os_volume_capacity").(int),
		})

		volStateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    VolumeSyncStatusRefreshFunc(service, osVolumeID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err := volStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", osVolumeID, ecloudservice.SyncStatusComplete, err)
		}
	}

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

			attachVolumeRequest := ecloudservice.AttachVolumeRequest{
				InstanceID: d.Id(),
			}
			log.Printf("[DEBUG] Created AttachVolumeRequest: %+v", attachVolumeRequest)

			log.Print("[INFO] Attaching volume")
			err := service.AttachVolume(volumeID, attachVolumeRequest)
			if err != nil {
				return fmt.Errorf("Error attaching volume: %s", err)
			}

			volStateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.SyncStatusComplete.String()},
				Refresh:    VolumeSyncStatusRefreshFunc(service, volumeID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = volStateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", volumeID, ecloudservice.SyncStatusComplete, err)
			}
		}

		//diff old against new for detach
		for _, id := range oldIDs {
			volumeID := id.(string)
			if rawVolumeExistsById(newIDs, volumeID) {
				continue
			}

			detachVolumeRequest := ecloudservice.DetachVolumeRequest{
				InstanceID: d.Id(),
			}
			log.Printf("[DEBUG] Created DetachVolumeRequest: %+v", detachVolumeRequest)

			log.Print("[INFO] Detaching volume")
			err := service.DetachVolume(volumeID, detachVolumeRequest)
			if err != nil {
				return fmt.Errorf("Error detaching volume: %s", err)
			}

			volStateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.SyncStatusComplete.String()},
				Refresh:    VolumeSyncStatusRefreshFunc(service, volumeID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = volStateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", volumeID, ecloudservice.SyncStatusComplete, err)
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

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    InstanceSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
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
