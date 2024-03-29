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

func resourceHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostGroupCreate,
		ReadContext:   resourceHostGroupRead,
		UpdateContext: resourceHostGroupUpdate,
		DeleteContext: resourceHostGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_spec_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"windows_enabled": {
				Type:     schema.TypeBool,
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
				Computed: true,
			},
		},
	}
}

func resourceHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateHostGroupRequest{
		VPCID:              d.Get("vpc_id").(string),
		HostSpecID:         d.Get("host_spec_id").(string),
		WindowsEnabled:     d.Get("windows_enabled").(bool),
		Name:               d.Get("name").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateHostGroupRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Host Group")
	task, err := service.CreateHostGroup(createReq)
	if err != nil {
		return diag.Errorf("Error creating host group: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for host group with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceHostGroupRead(ctx, d, meta)
}

func resourceHostGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving host group", map[string]interface{}{
		"id": d.Id(),
	})
	hostGroup, err := service.GetHostGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostGroupNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", hostGroup.VPCID)
	d.Set("host_spec_id", hostGroup.HostSpecID)
	d.Set("availability_zone_id", hostGroup.AvailabilityZoneID)
	d.Set("windows_enabled", hostGroup.WindowsEnabled)
	d.Set("name", hostGroup.Name)

	return nil
}

func resourceHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating host group", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchHostGroupRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchHostGroup(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating host group with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for host group with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceHostGroupRead(ctx, d, meta)
}

func resourceHostGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing host group", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteHostGroup(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostGroupNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing host group with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for host group with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
