package ecloud

import (
	"fmt"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

func resourceNIC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNICCreate,
		ReadContext:   resourceNICRead,
		UpdateContext: resourceNICUpdate,
		DeleteContext: resourceNICDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNICCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateNICRequest{
		Name:       d.Get("name").(string),
		InstanceID: d.Get("instance_id").(string),
		NetworkID:  d.Get("network_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateNICRequest: %+v", createReq))

	tflog.Info(ctx, "Creating NIC")
	taskRef, err := service.CreateNIC(createReq)
	if err != nil {
		return diag.Errorf("Error creating NIC: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for nic with ID [%s] to be created: %s", d.Id(), err)
	}

	d.SetId(taskRef.ResourceID)

	return resourceNICRead(ctx, d, meta)
}

func resourceNICRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving NIC", map[string]interface{}{
		"id": d.Id(),
	})
	nic, err := service.GetNIC(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NICNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("id", nic.ID)
	d.Set("name", nic.Name)
	d.Set("instance_id", nic.InstanceID)
	d.Set("network_id", nic.NetworkID)
	d.Set("mac_address", nic.MACAddress)

	if nic.IPAddress != "" {
		d.Set("ip_address", nic.IPAddress)
	}

	return nil
}

func resourceNICUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating NIC", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchNICRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchNIC(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating NIC with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for NIC with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceNICRead(ctx, d, meta)
}

func resourceNICDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	taskID, err := service.DeleteNIC(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NICNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing NIC with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for NIC with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
