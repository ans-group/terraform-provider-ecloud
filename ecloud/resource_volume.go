package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					iops := []int{300, 600, 1200, 2500, 5000}
					intInSlice := func(slice []int, value int) bool {
						for _, s := range slice {
							if s == value {
								return true
							}
						}
						return false
					}

					if !intInSlice(iops, v) {
						errs = append(errs, fmt.Errorf("%q must be a valid IOPS value [300, 600, 1200, 2500, 5000], got: %d", key, v))
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

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		createReq.IsShared = true
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateVolumeRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Volume")
	taskRef, err := service.CreateVolume(createReq)
	if err != nil {
		return diag.Errorf("Error creating volume: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for volume with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving volume", map[string]interface{}{
		"id": d.Id(),
	})
	volume, err := service.GetVolume(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
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

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		tflog.Info(ctx, "Updating volume", map[string]interface{}{
			"id": d.Id(),
		})
		task, err := service.PatchVolume(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating volume with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for volume with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}
	return resourceVolumeRead(ctx, d, meta)
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing volume", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVolume(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing volume with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for volume with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
