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

func resourceNICIPAddressBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNICIPAddressBindingCreate,
		ReadContext:   resourceNICIPAddressBindingRead,
		DeleteContext: resourceNICIPAddressBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"nic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceNICIPAddressBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	bindReq := ecloudservice.AssignIPAddressRequest{
		IPAddressID: d.Get("ip_address_id").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created AssignIPAddressRequest: %+v", bindReq))

	tflog.Info(ctx, "Assigning NIC IP address")
	taskID, err := service.AssignNICIPAddress(d.Get("nic_id").(string), bindReq)
	if err != nil {
		return diag.Errorf("Error assigning NIC IP address: %s", err)
	}

	d.SetId(getID(d.Get("nic_id").(string), d.Get("ip_address_id").(string)))

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for IP address with ID [%s] to be bound to NIC: %s", d.Id(), err)
	}

	return resourceNICIPAddressBindingRead(ctx, d, meta)
}

func resourceNICIPAddressBindingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	nicID := d.Get("nic_id").(string)
	ipAddressID := d.Get("ip_address_id").(string)

	ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
		*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
	))
	if err != nil {
		return diag.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
	}

	retrievedIPAddressID := ""
	if len(ipAddresses) == 1 {
		retrievedIPAddressID = ipAddresses[0].ID
	}

	d.Set("ip_address_id", retrievedIPAddressID)

	return nil
}

func resourceNICIPAddressBindingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	nicID := d.Get("nic_id").(string)
	ipAddressID := d.Get("ip_address_id").(string)

	ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
		*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
	))
	if err != nil {
		return diag.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
	}

	if len(ipAddresses) < 1 {
		return nil
	}

	tflog.Info(ctx, "Unassigning IP address from NIC", map[string]interface{}{
		"ip_address_id": ipAddressID,
		"nic_id":        nicID,
	})
	taskID, err := service.UnassignNICIPAddress(nicID, ipAddressID)
	if err != nil {
		return diag.Errorf("Error unassigning IP address [%s] from NIC [%s]: %s", ipAddressID, nicID, err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for IP address with ID [%s] to be unassigned from NIC [%s]: %s", ipAddressID, nicID, err)
	}

	return nil
}

func getID(nicID string, ipAddressID string) string {
	return fmt.Sprintf("%s.%s", nicID, ipAddressID)
}
