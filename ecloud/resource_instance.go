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
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
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

	d.Set("vpc_id", instance.VPCID)
	d.Set("name", instance.Name)
	d.Set("image_id", instance.ImageID)
	d.Set("vcpu_cores", instance.VCPUCores)
	d.Set("ram_capacity", instance.RAMCapacity)
	d.Set("volume_capacity", instance.VolumeCapacity)
	d.Set("locked", instance.Locked)
	d.Set("backup_enabled", instance.BackupEnabled)

	if d.Get("volume_id").(string) == "" {
		volumes, err := service.GetInstanceVolumes(d.Id(), connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Failed to retrieve instance volumes: %w", err)
		}

		if len(volumes) != 1 {
			return fmt.Errorf("Unexpected number of volumes (%d), expected 1", len(volumes))
		}

		d.Set("volume_id", volumes[0].ID)
	}

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

	if d.HasChange("volume_capacity") {
		log.Printf("[INFO] Updating volume with ID [%s]", d.Get("volume_id").(string))
		service.PatchVolume(d.Get("volume_id").(string), ecloudservice.PatchVolumeRequest{
			Capacity: d.Get("volume_capacity").(int),
		})

		// TODO: wait for volume sync
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
