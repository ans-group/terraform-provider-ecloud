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

func resourceVPNGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPNGatewayCreate,
		ReadContext:   resourceVPNGatewayRead,
		UpdateContext: resourceVPNGatewayUpdate,
		DeleteContext: resourceVPNGatewayDelete,
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
				Computed: true,
			},
			"specification_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPNGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPNGatewayRequest{
		RouterID:        d.Get("router_id").(string),
		Name:            d.Get("name").(string),
		SpecificationID: d.Get("specification_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateVPNGatewayRequest: %+v", createReq))

	tflog.Info(ctx, "Creating VPN Gateway")
	taskRef, err := service.CreateVPNGateway(createReq)
	if err != nil {
		return diag.Errorf("Error creating VPN gateway: %s", err)
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
		return diag.Errorf("Error waiting for VPN gateway with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceVPNGatewayRead(ctx, d, meta)
}

func resourceVPNGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving VPN Gateway", map[string]interface{}{
		"id": d.Id(),
	})
	gateway, err := service.GetVPNGateway(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNGatewayNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("router_id", gateway.RouterID)
	d.Set("name", gateway.Name)
	d.Set("specification_id", gateway.SpecificationID)
	d.Set("fqdn", gateway.FQDN)

	return nil
}

func resourceVPNGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPNGatewayRequest{
			Name: d.Get("name").(string),
		}

		tflog.Info(ctx, "Updating VPN Gateway", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.PatchVPNGateway(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating VPN gateway with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for VPN gateway with ID [%s] to be updated: %s", d.Id(), err)
		}
	}

	return resourceVPNGatewayRead(ctx, d, meta)
}

func resourceVPNGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing VPN Gateway", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVPNGateway(d.Id())
	if err != nil {
		return diag.Errorf("Error removing VPN gateway with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for VPN gateway with ID [%s] to be removed: %s", d.Id(), err)
	}

	return nil
}
