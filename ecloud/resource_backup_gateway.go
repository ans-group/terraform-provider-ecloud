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

func resourceBackupGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupGatewayCreate,
		ReadContext:   resourceBackupGatewayRead,
		UpdateContext: resourceBackupGatewayUpdate,
		DeleteContext: resourceBackupGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_spec_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBackupGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateBackupGatewayRequest{
		VPCID:         d.Get("vpc_id").(string),
		Name:          d.Get("name").(string),
		GatewaySpecID: d.Get("gateway_spec_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateBackupGatewayRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Backup Gateway")
	taskRef, err := service.CreateBackupGateway(createReq)
	if err != nil {
		return diag.Errorf("Error creating backup gateway: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for backup gateway with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceBackupGatewayRead(ctx, d, meta)
}

func resourceBackupGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving backup gateway", map[string]interface{}{
		"id": d.Id(),
	})
	backupGateway, err := service.GetBackupGateway(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.BackupGatewayNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", backupGateway.VPCID)
	d.Set("name", backupGateway.Name)
	d.Set("availability_zone_id", backupGateway.AvailabilityZoneID)
	d.Set("gateway_spec_id", backupGateway.GatewaySpecID)

	return nil
}

func resourceBackupGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating backup gateway", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchBackupGatewayRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchBackupGateway(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating backup gateway with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for backup gateway with ID [%s] to be updated: %s", d.Id(), err)
		}
	}

	return resourceBackupGatewayRead(ctx, d, meta)
}

func resourceBackupGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing backup gateway", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteBackupGateway(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.BackupGatewayNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing backup gateway with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for backup gateway with ID [%s] to be removed: %s", d.Id(), err)
	}

	return nil
}
