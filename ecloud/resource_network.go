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

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet": {
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

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateNetworkRequest{
		RouterID: d.Get("router_id").(string),
		Subnet:   d.Get("subnet").(string),
		Name:     d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateNetworkRequest: %+v", createReq))

	tflog.Info(ctx, "Creating network")
	networkID, err := service.CreateNetwork(createReq)
	if err != nil {
		return diag.Errorf("Error creating network: %s", err)
	}

	d.SetId(networkID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    NetworkSyncStatusRefreshFunc(service, networkID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for network with ID [%s] to return sync status of [%s]: %s", networkID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving network", map[string]interface{}{
		"id": d.Id(),
	})
	network, err := service.GetNetwork(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NetworkNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("router_id", network.RouterID)
	d.Set("subnet", network.Subnet)
	d.Set("name", network.Name)

	return nil
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchNetworkRequest{
			Name: d.Get("name").(string),
		}

		tflog.Info(ctx, "Updating network", map[string]interface{}{
			"id": d.Id(),
		})
		err := service.PatchNetwork(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating network with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    NetworkSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for network with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing network", map[string]interface{}{
		"id": d.Id(),
	})
	err := service.DeleteNetwork(d.Id())
	if err != nil {
		return diag.Errorf("Error removing network with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    NetworkSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for network with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// NetworkSyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func NetworkSyncStatusRefreshFunc(service ecloudservice.ECloudService, networkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		network, err := service.GetNetwork(networkID)
		if err != nil {
			if _, ok := err.(*ecloudservice.NetworkNotFoundError); ok {
				return network, "Deleted", nil
			}
			return nil, "", err
		}

		if network.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update network - review logs")
		}

		return network, network.Sync.Status.String(), nil
	}
}
