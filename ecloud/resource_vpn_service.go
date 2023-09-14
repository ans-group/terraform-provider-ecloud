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

func resourceVPNService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPNServiceCreate,
		ReadContext:   resourceVPNServiceRead,
		UpdateContext: resourceVPNServiceUpdate,
		DeleteContext: resourceVPNServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
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

func resourceVPNServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPNServiceRequest{
		RouterID: d.Get("router_id").(string),
		Name:     d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateVPNServiceRequest: %+v", createReq))

	tflog.Info(ctx, "Creating VPN service")
	taskRef, err := service.CreateVPNService(createReq)
	if err != nil {
		return diag.Errorf("Error creating VPN service: %s", err)
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
		return diag.Errorf("Error waiting for VPN service with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
	}

	return resourceVPNServiceRead(ctx, d, meta)
}

func resourceVPNServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving VPNService", map[string]interface{}{
		"id": d.Id(),
	})
	vpc, err := service.GetVPNService(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNServiceNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("router_id", vpc.RouterID)
	d.Set("name", vpc.Name)

	return nil
}

func resourceVPNServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPNServiceRequest{
			Name: d.Get("name").(string),
		}

		tflog.Info(ctx, "Updating VPNService", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.PatchVPNService(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating VPNService with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for VPN service with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceVPNServiceRead(ctx, d, meta)
}

func resourceVPNServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing VPNService", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVPNService(d.Id())
	if err != nil {
		return diag.Errorf("Error VPNService with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for VPNService with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
