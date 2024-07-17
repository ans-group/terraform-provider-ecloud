package ecloud

import (
	"context"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVolumeGroupInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeGroupInstanceCreate,
		ReadContext:   resourceVolumeGroupInstanceRead,
		DeleteContext: resourceVolumeGroupInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"volume_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVolumeGroupInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	instanceID := d.Get("instance_id").(string)
	volumeGroupID := d.Get("volume_group_id").(string)

	patchReq := ecloudservice.PatchInstanceRequest{
		VolumeGroupID: ptr.String(d.Get("volume_group_id").(string)),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created PatchInstanceRequest: %+v", patchReq))

	tflog.Info(ctx, "Updating instance", map[string]interface{}{
		"id": instanceID,
	})
	err := service.PatchInstance(instanceID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating instance with ID [%s]: %s", instanceID, err)
	}

	d.SetId(fmt.Sprintf("%s.%s", instanceID, volumeGroupID))

	_, err = waitForResourceState(
		ctx,
		ecloudservice.SyncStatusComplete.String(),
		InstanceSyncStatusRefreshFunc(service, instanceID),
		d.Timeout(schema.TimeoutUpdate),
	)
	if err != nil {
		return diag.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", instanceID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceVolumeGroupInstanceRead(ctx, d, meta)
}

func resourceVolumeGroupInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	instanceID := d.Get("instance_id").(string)

	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return diag.Errorf("Failed to retrieve instance: %s", err)
	}

	d.Set("volume_group_id", instance.VolumeGroupID)

	return nil
}

func resourceVolumeGroupInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	instanceID := d.Get("instance_id").(string)

	patchReq := ecloudservice.PatchInstanceRequest{
		VolumeGroupID: ptr.String(""),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created PatchInstanceRequest: %+v", patchReq))

	tflog.Info(ctx, "Updating instance", map[string]interface{}{
		"id": instanceID,
	})
	err := service.PatchInstance(instanceID, patchReq)
	if err != nil {
		return diag.Errorf("Error updating instance with ID [%s]: %s", instanceID, err)
	}

	_, err = waitForResourceState(
		ctx,
		ecloudservice.SyncStatusComplete.String(),
		InstanceSyncStatusRefreshFunc(service, instanceID),
		d.Timeout(schema.TimeoutUpdate),
	)
	if err != nil {
		return diag.Errorf("Error waiting for instance with ID [%s] to return sync status of [%s]: %s", instanceID, ecloudservice.SyncStatusComplete, err)
	}

	return nil
}
