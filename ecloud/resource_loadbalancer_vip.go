package ecloud

import (
	"context"
	"fmt"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLoadBalancerVip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerVipCreate,
		ReadContext:   resourceLoadBalancerVipRead,
		UpdateContext: resourceLoadBalancerVipUpdate,
		DeleteContext: resourceLoadBalancerVipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allocate_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLoadBalancerVipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVIPRequest{
		LoadBalancerID:     d.Get("load_balancer_id").(string),
		AllocateFloatingIP: d.Get("allocate_floating_ip").(bool),
		Name:               d.Get("name").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateVIPRequest: %+v", createReq))

	tflog.Info(ctx, "Creating LoadBalancer VIP")
	taskRef, err := service.CreateVIP(createReq)
	if err != nil {
		return diag.Errorf("Error creating loadbalancer vip: %s", err)
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
		return diag.Errorf("Error waiting for loadbalancer vip with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceLoadBalancerVipRead(ctx, d, meta)
}

func resourceLoadBalancerVipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving loadbalancer VIP", map[string]interface{}{
		"id": d.Id(),
	})
	lbVip, err := service.GetVIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", lbVip.Name)
	d.Set("load_balancer_id", lbVip.LoadBalancerID)

	if d.Get("floating_ip_id").(string) == "" && d.Get("allocate_floating_ip").(bool) {
		// we need to use the IP ID from the vip
		// and filter for the floating IP ID
		params := connection.APIRequestParameters{}
		params.WithFilter(*connection.NewAPIRequestFiltering("resource_id", connection.EQOperator, []string{lbVip.IPAddressID}))

		fips, err := service.GetFloatingIPs(params)
		if err != nil {
			return diag.Errorf("Failed to retrieve floating IPs: %s", err)
		}

		if len(fips) != 1 {
			return diag.Errorf("Unexpected number of floating IPs allocated to VIP")
		}

		d.Set("floating_ip_id", fips[0].ID)
	}

	return nil
}

func resourceLoadBalancerVipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating loadbalancer VIP", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchVIPRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchVIP(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating loadbalancer vip with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for loadbalancer vip with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceLoadBalancerVipRead(ctx, d, meta)
}

func resourceLoadBalancerVipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing loadbalancer VIP", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteVIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VIPNotFoundError:
			tflog.Debug(ctx, "Loadbalancer VIP not found, continuing", map[string]interface{}{
				"id": d.Id(),
			})
		default:
			return diag.Errorf("Error removing loadbalancer vip with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for loadbalancer vip with ID [%s] to be deleted: %s", d.Id(), err)
	}

	// remove floating ip if set
	if len(d.Get("floating_ip_id").(string)) > 1 {
		fip := d.Get("floating_ip_id").(string)

		tflog.Debug(ctx, "Removing floating IP", map[string]interface{}{
			"fip_id": fip,
		})

		taskID, err = service.DeleteFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				tflog.Info(ctx, "Floating IP not found, skipping delete", map[string]interface{}{
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
