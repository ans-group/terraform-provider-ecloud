package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNICIPAddressBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceNICIPAddressBindingCreate,
		Read:   resourceNICIPAddressBindingRead,
		Delete: resourceNICIPAddressBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceNICIPAddressBindingCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	bindReq := ecloudservice.AssignIPAddressRequest{
		IPAddressID: d.Get("ip_address_id").(string),
	}
	log.Printf("[DEBUG] Created AssignIPAddressRequest: %+v", bindReq)

	log.Print("[INFO] Assigning NIC IP address")
	taskID, err := service.AssignNICIPAddress(d.Get("nic_id").(string), bindReq)
	if err != nil {
		return fmt.Errorf("Error assigning NIC IP address: %s", err)
	}

	d.SetId(getID(d.Get("nic_id").(string), d.Get("ip_address_id").(string)))

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for IP address with ID [%s] to be bound to NIC: %s", d.Id(), err)
	}

	return resourceNICIPAddressBindingRead(d, meta)
}

func resourceNICIPAddressBindingRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	nicID := d.Get("nic_id").(string)
	ipAddressID := d.Get("ip_address_id").(string)

	ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
		*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
	))
	if err != nil {
		return fmt.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
	}

	retrievedIPAddressID := ""
	if len(ipAddresses) == 1 {
		retrievedIPAddressID = ipAddresses[0].ID
	}

	d.Set("ip_address_id", retrievedIPAddressID)

	return nil
}

func resourceNICIPAddressBindingDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	nicID := d.Get("nic_id").(string)
	ipAddressID := d.Get("ip_address_id").(string)

	ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
		*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
	))
	if err != nil {
		return fmt.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
	}

	if len(ipAddresses) < 1 {
		return nil
	}

	log.Printf("[INFO] Unassigning IP address [%s] from NIC [%s]", ipAddressID, nicID)
	taskID, err := service.UnassignNICIPAddress(nicID, ipAddressID)
	if err != nil {
		return fmt.Errorf("Error unassigning IP address [%s] from NIC [%s]: %s", ipAddressID, nicID, err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for IP address with ID [%s] to be unassigned from NIC [%s]: %s", ipAddressID, nicID, err)
	}

	return nil
}

func getID(nicID string, ipAddressID string) string {
	return fmt.Sprintf("%s.%s", nicID, ipAddressID)
}
