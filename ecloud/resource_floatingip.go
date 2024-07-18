package ecloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFloatingIPCreate,
		ReadContext:   resourceFloatingIPRead,
		UpdateContext: resourceFloatingIPUpdate,
		DeleteContext: resourceFloatingIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					id := val.(string)
					fipAssignableResources := []string{"nic-", "ip-", "rtr-"}

					prefixInSlice := func(slice []string, value string) bool {
						for _, s := range slice {
							if strings.HasPrefix(value, s) {
								return true
							}
						}
						return false
					}

					if !prefixInSlice(fipAssignableResources, id) {
						errs = append(errs, fmt.Errorf("%q must be a valid resource that supports floating ip assignment. got: %s", key, id))
					}
					return
				},
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFloatingIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateFloatingIPRequest{
		Name:               d.Get("name").(string),
		VPCID:              d.Get("vpc_id").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Created CreateFloatingIPRequest: %+v", createReq))

	tflog.Info(ctx, "Creating Floating IP")
	taskRef, err := service.CreateFloatingIP(createReq)
	if err != nil {
		return diag.Errorf("Error creating floating IP: %s", err)
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
		return diag.Errorf("Error waiting for floating IP with ID [%s] to be created: %s", d.Id(), err)
	}

	if r, ok := d.GetOk("resource_id"); ok {
		resourceID := r.(string)
		if strings.HasPrefix(resourceID, "nic-") {
			nicDHCPAddress, err := getNICDHCPAddress(service, resourceID)
			if err != nil {
				return diag.Errorf("Error retrieving DHCP IP address for NIC with ID [%s]: %s", resourceID, err)
			}
			resourceID = nicDHCPAddress.ID
		}

		tflog.Info(ctx, "Assigning floating IP", map[string]interface{}{
			"fip_id":          d.Id(),
			"target_resource": resourceID,
		})

		assignFipReq := ecloudservice.AssignFloatingIPRequest{
			ResourceID: r.(string),
		}
		tflog.Debug(ctx, fmt.Sprintf("Created AssignFloatingIPRequest: %+v", assignFipReq))

		taskID, err := service.AssignFloatingIP(d.Id(), assignFipReq)
		if err != nil {
			return diag.Errorf("Error assigning floating IP: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for floating IP with ID [%s] to be assigned: %s", d.Id(), err)
		}
	}
	return resourceFloatingIPRead(ctx, d, meta)
}

func resourceFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving floating IP", map[string]interface{}{
		"id": d.Id(),
	})
	fip, err := service.GetFloatingIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FloatingIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("vpc_id", fip.VPCID)
	d.Set("name", fip.Name)
	d.Set("ip_address", fip.IPAddress)
	d.Set("availability_zone_id", fip.AvailabilityZoneID)

	return nil
}

func resourceFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating floating IP", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchFloatingIPRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchFloatingIP(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating floating ip with ID [%s]: %s", d.Id(), err)
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
			return diag.Errorf("Error waiting for floating ip with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("resource_id") {
		tflog.Info(ctx, "Updating floating IP", map[string]interface{}{
			"id": d.Id(),
		})

		oldVal, newVal := d.GetChange("resource_id")
		assign := true

		// Handle scenario where user is updating resource from NIC to corresponding DHCP address for that NIC.
		// Here we check whether the provided IP address matches the DHCP IP address of the previously defined
		// NIC, and if so, skip unassign/assign
		if strings.HasPrefix(oldVal.(string), "nic-") && strings.HasPrefix(newVal.(string), "ip-") {
			nicID := oldVal.(string)
			nicDHCPAddress, err := getNICDHCPAddress(service, nicID)
			if err != nil {
				return diag.Errorf("Error retrieving DHCP IP address for NIC with ID [%s]: %s", nicID, err)
			}

			if nicDHCPAddress.ID == newVal.(string) {
				assign = false
			}
		}

		// Handle scenario where user is updating resource from IP to NIC which corresponds to the DHCP address for that NIC.
		// Here we check whether the provided NIC has a DHCP IP address which matches the ID of the previously defined IP address,
		// and if so, skip unassign/assign
		if strings.HasPrefix(oldVal.(string), "ip-") && strings.HasPrefix(newVal.(string), "nic-") {
			nicID := newVal.(string)
			nicDHCPAddress, err := getNICDHCPAddress(service, nicID)
			if err != nil {
				return diag.Errorf("Error retrieving DHCP IP address for NIC with ID [%s]: %s", nicID, err)
			}

			if nicDHCPAddress.ID == oldVal.(string) {
				assign = false
			}
		}

		// if oldVal wasn't empty then floating ip needs unassigned first
		if assign && oldVal.(string) != "" {
			tflog.Debug(ctx, "Unassigning floating IP", map[string]interface{}{
				"id": d.Id(),
			})
			taskID, err := service.UnassignFloatingIP(d.Id())
			if err != nil {
				return diag.Errorf("Error unassigning floating ip with ID [%s]: %s", d.Id(), err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      3 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %s", d.Id(), err)
			}
		}

		// Assign floating ip to new instance value if set
		if assign && len(newVal.(string)) > 1 {
			tflog.Info(ctx, "Assigning floating IP", map[string]interface{}{
				"fip_id":          d.Id(),
				"target_resource": newVal.(string),
			})

			assignFipReq := ecloudservice.AssignFloatingIPRequest{
				ResourceID: newVal.(string),
			}
			tflog.Debug(ctx, fmt.Sprintf("Created AssignFloatingIPRequest: %+v", assignFipReq))

			taskID, err := service.AssignFloatingIP(d.Id(), assignFipReq)
			if err != nil {
				return diag.Errorf("Error assigning floating IP: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      3 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for floating IP with ID [%s] to be assigned: %s", d.Id(), err)
			}
		}
	}
	return resourceFloatingIPRead(ctx, d, meta)
}

func resourceFloatingIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	// first check if floating ip is assigned
	if _, ok := d.GetOk("resource_id"); ok {
		taskID, err := service.UnassignFloatingIP(d.Id())
		if err != nil {
			return diag.Errorf("Error unassigning floating ip with ID [%s]: %s", d.Id(), err)
		}

		unassignStateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = unassignStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %s", d.Id(), err)
		}
	}

	// Once unassigned - remove the floating ip resource
	tflog.Info(ctx, "Removing floating IP", map[string]interface{}{
		"id": d.Id(),
	})
	taskID, err := service.DeleteFloatingIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FloatingIPNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing floating ip with ID [%s]: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for floating ip with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

func getNICDHCPAddress(service ecloudservice.ECloudService, nicID string) (ecloudservice.IPAddress, error) {
	filter := connection.NewAPIRequestParameters().
		WithFilter(connection.APIRequestFiltering{
			Property: "type",
			Operator: connection.EQOperator,
			Value:    []string{"dhcp"},
		})
	ipAddresses, err := service.GetNICIPAddresses(nicID, *filter)
	if err != nil {
		return ecloudservice.IPAddress{}, fmt.Errorf("Error retrieving IP addresses for NIC with ID [%s]: %s", nicID, err)
	}

	if len(ipAddresses) != 1 {
		return ecloudservice.IPAddress{}, fmt.Errorf("Unexpected number of DHCP IP addresses [%d] for NIC with ID [%s], expected 1", len(ipAddresses), nicID)
	}

	return ipAddresses[0], nil
}
