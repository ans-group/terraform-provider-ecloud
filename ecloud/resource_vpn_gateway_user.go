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

func resourceVPNGatewayUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPNGatewayUserCreate,
		ReadContext:   resourceVPNGatewayUserRead,
		UpdateContext: resourceVPNGatewayUserUpdate,
		DeleteContext: resourceVPNGatewayUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				// XXX: Once password change functionality is properly implement this can be removed,
				// until then we need to make sure that we're creating a new user if a password needs
				// updating.
				// See ADO#36490
				ForceNew: true,
			},
		},
	}
}

func resourceVPNGatewayUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPNGatewayUserRequest{
		VPNGatewayID: d.Get("vpn_gateway_id").(string),
		Name:         d.Get("name").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateVPNGatewayUserRequest: %+v", createReq))

	tflog.Info(ctx, "Creating VPN Gateway User")
	taskRef, err := service.CreateVPNGatewayUser(createReq)
	if err != nil {
		return diag.Errorf("Error creating VPN gateway user: %s", err)
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
		return diag.Errorf("Error waiting for VPN gateway user with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceVPNGatewayUserRead(ctx, d, meta)
}

func resourceVPNGatewayUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving VPN Gateway User", map[string]interface{}{
		"id": d.Id(),
	})
	user, err := service.GetVPNGatewayUser(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNGatewayUserNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpn_gateway_id", user.VPNGatewayID)
	d.Set("name", user.Name)
	d.Set("username", user.Username)

	return nil
}

func resourceVPNGatewayUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	hasChange := false
	patchReq := ecloudservice.PatchVPNGatewayUserRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("password") {
		hasChange = true
		patchReq.Password = d.Get("password").(string)
	}

	if hasChange {
		tflog.Info(ctx, "Updating VPN Gateway User", map[string]interface{}{
			"id": d.Id(),
		})
		taskRef, err := service.PatchVPNGatewayUser(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating VPN gateway user with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for VPN gateway user with ID [%s] to be updated: %s", d.Id(), err)
		}
	}

	return resourceVPNGatewayUserRead(ctx, d, meta)
}

func resourceVPNGatewayUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing VPN Gateway User", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVPNGatewayUser(d.Id())
	if err != nil {
		return diag.Errorf("Error removing VPN gateway user with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for VPN gateway user with ID [%s] to be removed: %s", d.Id(), err)
	}

	return nil
}
