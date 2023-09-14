package ecloud

import (
	"context"
	"fmt"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageUpdate,
		DeleteContext: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	instanceID := d.Get("instance_id").(string)
	createReq := ecloudservice.CreateInstanceImageRequest{
		Name: d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateImageRequest: %+v", createReq))

	tflog.Info(ctx, "Creating image")
	taskRef, err := service.CreateInstanceImage(instanceID, createReq)
	if err != nil {
		return diag.Errorf("Error creating image: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceImageRead(ctx, d, meta)
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving image", map[string]interface{}{
		"id": d.Id(),
	})
	image, err := service.GetImage(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.ImageNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", image.VPCID)
	d.Set("name", image.Name)
	d.Set("availability_zone_id", image.AvailabilityZoneID)

	return nil
}

func resourceImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	hasChange := false
	patchReq := ecloudservice.UpdateImageRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if hasChange {
		tflog.Info(ctx, "Updating image", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.UpdateImage(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating image with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceImageRead(ctx, d, meta)
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing image", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteImage(d.Id())
	if err != nil {
		return diag.Errorf("Error removing image with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for image with ID [%s] to be deleted: %s", taskID, err)
	}

	return nil
}
