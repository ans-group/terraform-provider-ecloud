package ecloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ans-group/sdk-go/pkg/ptr"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCCreate,
		ReadContext:   resourceVPCRead,
		UpdateContext: resourceVPCUpdate,
		DeleteContext: resourceVPCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"advanced_networking": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPCRequest{
		RegionID:           d.Get("region_id").(string),
		Name:               d.Get("name").(string),
		ClientID:           d.Get("client_id").(int),
		AdvancedNetworking: ptr.Bool(d.Get("advanced_networking").(bool)),
	}
	log.Printf("[DEBUG] Created CreateVPCRequest: %+v", createReq)

	log.Print("[INFO] Creating VPC")
	vpcID, err := service.CreateVPC(createReq)
	if err != nil {
		return diag.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpcID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    VPCSyncStatusRefreshFunc(service, vpcID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for VPC with ID [%s] to return sync status of [%s]: %s", vpcID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceVPCRead(ctx, d, meta)
}

func resourceVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving VPC with ID [%s]", d.Id())
	vpc, err := service.GetVPC(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPCNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("region_id", vpc.RegionID)
	d.Set("name", vpc.Name)
	d.Set("advanced_networking", vpc.AdvancedNetworking)

	return nil
}

func resourceVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPCRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating VPC with ID [%s]", d.Id())
		err := service.PatchVPC(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating VPC with ID [%s]: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    VPCSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for VPC with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceVPCRead(ctx, d, meta)
}

func resourceVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing VPC with ID [%s]", d.Id())
	err := service.DeleteVPC(d.Id())
	if err != nil {
		return diag.Errorf("Error VPC with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    VPCSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for VPC with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

// VPCSyncStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func VPCSyncStatusRefreshFunc(service ecloudservice.ECloudService, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpc, err := service.GetVPC(vpcID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VPCNotFoundError); ok {
				return vpc, "Deleted", nil
			}
			return nil, "", err
		}

		if vpc.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update VPC - review logs")
		}

		return vpc, vpc.Sync.Status.String(), nil
	}
}
