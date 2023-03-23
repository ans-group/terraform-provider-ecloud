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

func resourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAddressCreate,
		Read:   resourceIPAddressRead,
		Update: resourceIPAddressUpdate,
		Delete: resourceIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateIPAddressRequest{
		NetworkID: d.Get("network_id").(string),
		Name:      d.Get("name").(string),
		IPAddress: connection.IPAddress(d.Get("ip_address").(string)),
	}
	log.Printf("[DEBUG] Created CreateIPAddressRequest: %+v", createReq)

	log.Print("[INFO] Creating IPAddress")
	task, err := service.CreateIPAddress(createReq)
	if err != nil {
		return fmt.Errorf("Error creating IP address: %s", err)
	}

	d.SetId(task.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for IP address with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceIPAddressRead(d, meta)
}

func resourceIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving IP address with ID [%s]", d.Id())
	ipAddress, err := service.GetIPAddress(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.IPAddressNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("network_id", ipAddress.NetworkID)
	d.Set("name", ipAddress.Name)
	d.Set("ip_address", ipAddress.IPAddress)

	return nil
}

func resourceIPAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating IP address with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchIPAddressRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchIPAddress(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating IP address with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, task.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for IP address with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceIPAddressRead(d, meta)
}

func resourceIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing IP address with ID [%s]", d.Id())
	taskID, err := service.DeleteIPAddress(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.IPAddressNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing IP address with ID [%s]: %s", d.Id(), err)
		}
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
		return fmt.Errorf("Error waiting for IP address with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
