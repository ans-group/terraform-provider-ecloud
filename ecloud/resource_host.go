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

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"host_group_id": {
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateHostRequest{
		HostGroupID: d.Get("host_group_id").(string),
		Name:        d.Get("name").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateHostRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Host Group")
	task, err := service.CreateHost(createReq)
	if err != nil {
		return diag.Errorf("Error creating host group: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for host group with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceHostRead(ctx, d, meta)
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving host group", map[string]interface{}{
		"id": d.Id(),
	})
	host, err := service.GetHost(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("host_group_id", host.HostGroupID)
	d.Set("name", host.Name)

	return nil
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating Host", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchHostRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchHost(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating host with ID [%s]: %s", d.Id(), err)
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

	return resourceHostRead(ctx, d, meta)
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing host group", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteHost(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing host with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for host with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
