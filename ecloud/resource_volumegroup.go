package ecloud

import (
	"context"
	"log"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVolumeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeGroupCreate,
		ReadContext:   resourceVolumeGroupRead,
		UpdateContext: resourceVolumeGroupUpdate,
		DeleteContext: resourceVolumeGroupDelete,
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
		},
	}
}

func resourceVolumeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVolumeGroupRequest{
		VPCID:              d.Get("vpc_id").(string),
		Name:               d.Get("name").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}

	log.Printf("[DEBUG] Created CreateVolumeGroupRequest: %+v", createReq)

	log.Print("[INFO] Creating VolumeGroup")
	taskRef, err := service.CreateVolumeGroup(createReq)
	if err != nil {
		return diag.Errorf("Error creating volumegroup: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for volumegroup with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceVolumeGroupRead(ctx, d, meta)
}

func resourceVolumeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving volume with ID [%s]", d.Id())
	volumegroup, err := service.GetVolumeGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeGroupNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", volumegroup.VPCID)
	d.Set("name", volumegroup.Name)
	d.Set("availability_zone_id", volumegroup.AvailabilityZoneID)

	return nil
}

func resourceVolumeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	patchReq := ecloudservice.PatchVolumeGroupRequest{}
	hasChange := false
	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[INFO] Updating volumegroup with ID [%s]", d.Id())
		task, err := service.PatchVolumeGroup(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating volumegroup with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for volumegroup with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}
	return resourceVolumeGroupRead(ctx, d, meta)
}

func resourceVolumeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing volumegroup with ID [%s]", d.Id())
	taskID, err := service.DeleteVolumeGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VolumeGroupNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing volumegroup with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for volumegroup with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
