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

func resourceRouter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouterCreate,
		ReadContext:   resourceRouterRead,
		UpdateContext: resourceRouterUpdate,
		DeleteContext: resourceRouterDelete,
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
			"router_throughput_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateRouterRequest{
		VPCID:              d.Get("vpc_id").(string),
		Name:               d.Get("name").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
		RouterThroughputID: d.Get("router_throughput_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateRouterRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Router")
	routerID, err := service.CreateRouter(createReq)
	if err != nil {
		return diag.Errorf("Error creating router: %s", err)
	}

	d.SetId(routerID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    RouterSyncStatusRefreshFunc(service, routerID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for router with ID [%s] to return sync status of [%s]: %s", routerID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceRouterRead(ctx, d, meta)
}

func resourceRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving router", map[string]interface{}{
		"id": d.Id(),
	})
	router, err := service.GetRouter(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.RouterNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", router.VPCID)
	d.Set("name", router.Name)
	d.Set("availability_zone_id", router.AvailabilityZoneID)
	d.Set("router_throughput_id", router.RouterThroughputID)

	return nil
}

func resourceRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	hasChange := false
	patchReq := ecloudservice.PatchRouterRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("router_throughput_id") {
		hasChange = true
		patchReq.RouterThroughputID = d.Get("router_throughput_id").(string)
	}

	if hasChange {
		tflog.Info(ctx, "Updating router", map[string]interface{}{
			"id": d.Id(),
		})
		err := service.PatchRouter(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating router with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    RouterSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for router with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceRouterRead(ctx, d, meta)
}

func resourceRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing router", map[string]interface{}{
		"id": d.Id(),
	})
	err := service.DeleteRouter(d.Id())
	if err != nil {
		return diag.Errorf("Error removing router with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    RouterSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for router with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// RouterSyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func RouterSyncStatusRefreshFunc(service ecloudservice.ECloudService, routerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		router, err := service.GetRouter(routerID)
		if err != nil {
			if _, ok := err.(*ecloudservice.RouterNotFoundError); ok {
				return router, "Deleted", nil
			}
			return nil, "", err
		}

		if router.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update router - review logs")
		}

		return router, router.Sync.Status.String(), nil
	}
}
