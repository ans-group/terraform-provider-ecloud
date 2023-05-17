package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/terraform-provider-ecloud/pkg/lock"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			"image_data": {
				Type:     schema.TypeMap,
				Optional: true,
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
				ForceNew: true,
				Default:  false,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"requires_floating_ip"},
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if d.Get("requires_floating_ip").(bool) {
						return true
					}
					return oldValue == newValue && newValue == ""
				},
			},
			"requires_floating_ip": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				ConflictsWith: []string{"floating_ip_id"},
			},
			"data_volume_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
			"ssh_keypair_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				ForceNew: true,
			},
			"host_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_tier_id"},
			},
			"resource_tier_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"host_group_id"},
			},
			"volume_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
		ImageData:          expandCreateInstanceRequestImageData(d.Get("image_data").(map[string]interface{})),
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
		IsEncrypted:        d.Get("encrypted").(bool),
		HostGroupID:        d.Get("host_group_id").(string),
		ResourceTierID:     d.Get("resource_tier_id").(string),
		CustomIPAddress:    connection.IPAddress(d.Get("ip_address").(string)),
		SSHKeyPairIDs:      expandSshKeyPairIds(d.Get("ssh_keypair_ids").(*schema.Set).List()),
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

	//attach data volumes
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
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to be attached: %s", volumeID, err)
			}
		}
	}

	//handle volume group if defined
	if volumeGroupID, ok := d.GetOk("volume_group_id"); ok {
		patchReq := ecloudservice.PatchInstanceRequest{
			VolumeGroupID: ptr.String(volumeGroupID.(string)),
		}
		log.Printf("[DEBUG] Created PatchInstanceRequest: %+v", patchReq)

		err := service.PatchInstance(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error attaching volume: %s", err)
		}

		_, err = waitForResourceState(
			ecloudservice.SyncStatusComplete.String(),
			InstanceSyncStatusRefreshFunc(service, d.Id()),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for volumegroup with ID [%s] to be attached: %s", volumeGroupID, err)
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
	d.Set("host_group_id", instance.HostGroupID)
	d.Set("resource_tier_id", instance.ResourceTierID)
	d.Set("volume_capacity", osVolume[0].Capacity)
	d.Set("volume_iops", osVolume[0].IOPS)
	d.Set("volume_group_id", instance.VolumeGroupID)
	d.Set("encrypted", instance.IsEncrypted)

	if d.Get("nic_id").(string) == "" {
		nics, err := service.GetInstanceNICs(d.Id(), connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Failed to retrieve instance nics: %w", err)
		}

		if len(nics) > 1 {
			return fmt.Errorf("Unexpected number of instance nics. Unable to lookup floating ip")
		}

		d.Set("nic_id", nics[0].ID)
	}

	if d.Get("requires_floating_ip").(bool) {
		//we need to retrieve the instance nic to find the associated floating ip
		fips, err := service.GetInstanceFloatingIPs(d.Id(), connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Failed to retrieve floating IPs: %w", err)
		}

		if len(fips) != 1 {
			return fmt.Errorf("Unexpected number of floating IPs assigned to instance")
		}

		d.Set("floating_ip_id", fips[0].ID)
	}

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
	if d.HasChange("volume_group_id") {
		hasChange = true
		patchReq.VolumeGroupID = ptr.String(d.Get("volume_group_id").(string))
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

	if d.HasChange("floating_ip_id") && !d.Get("requires_floating_ip").(bool) {

		oldVal, newVal := d.GetChange("floating_ip_id")
		oldFip := oldVal.(string)
		newFip := newVal.(string)

		log.Printf("[DEBUG] ##Floating IP change detected. oldValue: [%s], newValue [%s]", oldFip, newFip)

		if len(newFip) < 1 && oldFip != "" {
			// lock the fip
			if len(oldFip) < 1 {
				return fmt.Errorf("invalid floating ip ID: %s", oldFip)
			}

			unlock := lock.LockResource(oldFip)
			defer unlock()

			log.Printf("[DEBUG] Unassigning floating IP with ID [%s]", oldFip)

			//unassign floating ip but don't delete as it may be managed by another resource
			taskID, err := service.UnassignFloatingIP(oldFip)
			if err != nil {
				switch err.(type) {
				case *ecloudservice.FloatingIPNotFoundError:
					log.Printf("[DEBUG] Floating IP with ID [%s] not found. Skipping unassign.", oldFip)
				default:
					return fmt.Errorf("Error unassigning floating ip with ID [%s]: %s", oldFip, err)
				}
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutDelete),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %w", oldFip, err)
			}

			//unset floating ip
			d.Set("floating_ip_id", "")
		}

		if oldFip == "" && newFip != "" {
			// lock the fip
			if len(newFip) < 1 {
				return fmt.Errorf("invalid floating ip ID: %s", newFip)
			}

			unlock := lock.LockResource(newFip)
			defer unlock()

			log.Printf("[DEBUG] Assigning floating ip with ID [%s]", newFip)

			assignFipReq := ecloudservice.AssignFloatingIPRequest{
				ResourceID: d.Get("nic_id").(string),
			}
			log.Printf("[DEBUG] Created AssignFloatingIPRequest: %+v", assignFipReq)

			taskID, err := service.AssignFloatingIP(newFip, assignFipReq)
			if err != nil {
				return fmt.Errorf("Error assigning floating IP: %s", err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for floating IP with ID [%s] to be assigned: %s", newFip, err)
			}
		}
	}

	//manage volume capacity
	if d.HasChange("volume_capacity") {
		osVolumeID := d.Get("volume_id").(string)
		log.Printf("[INFO] Updating volume with ID [%s]", osVolumeID)
		task, err := service.PatchVolume(osVolumeID, ecloudservice.PatchVolumeRequest{
			Capacity: d.Get("volume_capacity").(int),
		})
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", osVolumeID, err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, task.TaskID),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", osVolumeID, ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("volume_iops") {
		osVolumeID := d.Get("volume_id").(string)
		log.Printf("[INFO] Updating volume with ID [%s]", osVolumeID)
		task, err := service.PatchVolume(osVolumeID, ecloudservice.PatchVolumeRequest{
			IOPS: d.Get("volume_iops").(int),
		})
		if err != nil {
			return fmt.Errorf("Error updating volume with ID [%s]: %w", osVolumeID, err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, task.TaskID),
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
				TaskStatusRefreshFunc(service, taskID),
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
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", volumeID, ecloudservice.TaskStatusComplete, err)
			}
		}
	}

	// Handle host groups
	if d.HasChange("host_group_id") {
		hostGroupID := d.Get("host_group_id").(string)
		migrateReq := ecloud.MigrateInstanceRequest{
			HostGroupID: hostGroupID,
		}

		log.Printf("[INFO] Migrating instance [%s] to host group [%s]", d.Id(), hostGroupID)
		taskID, err := service.MigrateInstance(d.Id(), migrateReq)
		if err != nil {
			return fmt.Errorf("Error migrating instance: %s", err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
		}
	}

	// handle resource tier migrations if hostgroup id isn't populated
	if d.HasChange("resource_tier_id") {
		resourceTierID := d.Get("resource_tier_id").(string)

		migrateReq := ecloud.MigrateInstanceRequest{
			ResourceTierID: resourceTierID,
		}

		log.Printf("[INFO] Migrating instance [%s] to resource tier [%s]", d.Id(), resourceTierID)
		taskID, err := service.MigrateInstance(d.Id(), migrateReq)
		if err != nil {
			return fmt.Errorf("Error migrating instance: %s", err)
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutUpdate),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("encrypted") {
		isEncrypted := d.Get("encrypted").(bool)
		log.Printf("[INFO] Updating instance encryption status to [%s]", isEncrypted)

		if isEncrypted {
			taskID, err := service.EncryptInstance(d.Id())
			if err != nil {
				return fmt.Errorf("Error encrypting instance [%s]: %w", d.Id(), err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
			}
		} else {
			taskID, err := service.DecryptInstance(d.Id())
			if err != nil {
				return fmt.Errorf("Error decrypting instance [%s]: %w", d.Id(), err)
			}

			_, err = waitForResourceState(
				ecloudservice.TaskStatusComplete.String(),
				TaskStatusRefreshFunc(service, taskID),
				d.Timeout(schema.TimeoutUpdate),
			)
			if err != nil {
				return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskID, ecloudservice.TaskStatusComplete, err)
			}
		}
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	//remove floating ip if set
	if d.Get("requires_floating_ip").(bool) && len(d.Get("floating_ip_id").(string)) > 1 {
		fip := d.Get("floating_ip_id").(string)

		log.Printf("[DEBUG] Unassigning floating ip with ID [%s]", fip)

		taskID, err := service.UnassignFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				log.Printf("[DEBUG] Floating IP with ID [%s] not found. Skipping unassign.", fip)
			default:
				return fmt.Errorf("Error unassigning floating ip with ID [%s]: %s", fip, err)
			}
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutDelete),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %w", d.Id(), err)
		}

		log.Printf("[DEBUG] Removing floating ip with ID [%s]", fip)

		taskID, err = service.DeleteFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				log.Printf("[DEBUG] Floating IP with ID [%s] not found. Skipping delete.", fip)
			default:
				return fmt.Errorf("Error removing floating ip with ID [%s]: %s", fip, err)
			}
		}

		_, err = waitForResourceState(
			ecloudservice.TaskStatusComplete.String(),
			TaskStatusRefreshFunc(service, taskID),
			d.Timeout(schema.TimeoutDelete),
		)
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to be removed: %w", d.Id(), err)
		}
	}

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
