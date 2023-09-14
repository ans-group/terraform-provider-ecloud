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

func resourceVPNEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPNEndpointCreate,
		ReadContext:   resourceVPNEndpointRead,
		UpdateContext: resourceVPNEndpointUpdate,
		DeleteContext: resourceVPNEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpn_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"manage_floating_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceVPNEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	// if not populated then assume the provider should manage the fip
	if len(d.Get("floating_ip_id").(string)) < 1 {
		d.Set("manage_floating_ip", true)
	} else {
		d.Set("manage_floating_ip", false)
	}

	createReq := ecloudservice.CreateVPNEndpointRequest{
		VPNServiceID: d.Get("vpn_service_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		Name:         d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateVPNEndpointRequest: %+v", createReq))

	tflog.Info(ctx, "Creating VPN endpoint")
	taskRef, err := service.CreateVPNEndpoint(createReq)
	if err != nil {
		return diag.Errorf("Error creating VPN endpoint: %s", err)
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
		return diag.Errorf("Error waiting for VPN endpoint with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
	}

	return resourceVPNEndpointRead(ctx, d, meta)
}

func resourceVPNEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving VPNEndpoint", map[string]interface{}{
		"id": d.Id(),
	})
	vpc, err := service.GetVPNEndpoint(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNEndpointNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpn_service_id", vpc.VPNServiceID)
	d.Set("floating_ip_id", vpc.FloatingIPID)
	d.Set("name", vpc.Name)

	return nil
}

func resourceVPNEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPNEndpointRequest{
			Name: d.Get("name").(string),
		}

		tflog.Info(ctx, "Updating VPNEndpoint", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.PatchVPNEndpoint(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating VPNEndpoint with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for VPN endpoint with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceVPNEndpointRead(ctx, d, meta)
}

func resourceVPNEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing VPNEndpoint", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVPNEndpoint(d.Id())
	if err != nil {
		return diag.Errorf("Error VPNEndpoint with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for VPNEndpoint with ID [%s] to be deleted: %s", d.Id(), err)
	}

	// remove floating ip if set
	if d.Get("manage_floating_ip").(bool) {
		fip := d.Get("floating_ip_id").(string)

		tflog.Debug(ctx, "Removing floating IP", map[string]interface{}{
			"id": fip,
		})

		taskID, err = service.DeleteFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				tflog.Debug(ctx, "Floating IP not found, skipping delete", map[string]interface{}{
					"id": fip,
				})
			default:
				return diag.Errorf("Error removing floating ip with ID [%s]: %s", fip, err)
			}
		}

		stateConf = &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for floating ip with ID [%s] to be removed: %s", d.Id(), err)
		}
	}

	return nil
}
