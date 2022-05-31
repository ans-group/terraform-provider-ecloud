package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVPNSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPNSessionCreate,
		Read:   resourceVPNSessionRead,
		Update: resourceVPNSessionUpdate,
		Delete: resourceVPNSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpn_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_profile_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpn_endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_networks": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_networks": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"psk": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceVPNSessionCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPNSessionRequest{
		VPNServiceID:      d.Get("vpn_service_id").(string),
		VPNProfileGroupID: d.Get("vpn_profile_group_id").(string),
		VPNEndpointID:     d.Get("vpn_endpoint_id").(string),
		RemoteIP:          connection.IPAddress(d.Get("remote_ip").(string)),
		Name:              d.Get("name").(string),
		RemoteNetworks:    d.Get("remote_networks").(string),
		LocalNetworks:     d.Get("local_networks").(string),
	}
	log.Printf("[DEBUG] Created CreateVPNSessionRequest: %+v", createReq)

	log.Print("[INFO] Creating VPN session")
	taskRef, err := service.CreateVPNSession(createReq)
	if err != nil {
		return fmt.Errorf("Error creating VPN session: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPN session with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
	}

	if d.HasChange("psk") {
		updatePSKReq := ecloudservice.UpdateVPNSessionPreSharedKeyRequest{
			PSK: d.Get("psk").(string),
		}
		log.Printf("[DEBUG] Created UpdateVPNSessionPreSharedKeyRequest: %+v", updatePSKReq)

		taskRef, err := service.UpdateVPNSessionPreSharedKey(taskRef.ResourceID, updatePSKReq)
		if err != nil {
			return fmt.Errorf("Error creating VPN session pre-shared key: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for VPN session with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceVPNSessionRead(d, meta)
}

func resourceVPNSessionRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving VPN session with ID [%s]", d.Id())
	session, err := service.GetVPNSession(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPNSessionNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpn_service_id", session.VPNServiceID)
	d.Set("vpn_profile_group_id", session.VPNProfileGroupID)
	d.Set("vpn_endpoint_id", session.VPNEndpointID)
	d.Set("remote_ip", session.RemoteIP)
	d.Set("name", session.Name)
	d.Set("remote_networks", session.RemoteNetworks)
	d.Set("local_networks", session.LocalNetworks)

	log.Printf("[INFO] Retrieving VPN session pre-shared key with ID [%s]", d.Id())
	psk, err := service.GetVPNSessionPreSharedKey(d.Id())
	if err != nil {
		return err
	}
	d.Set("psk", psk.PSK)

	return nil
}

func resourceVPNSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	hasChange := false
	patchReq := ecloudservice.PatchVPNSessionRequest{}
	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}
	if d.HasChange("vpn_profile_group_id") {
		hasChange = true
		patchReq.VPNProfileGroupID = d.Get("vpn_profile_group_id").(string)
	}
	if d.HasChange("remote_ip") {
		hasChange = true
		patchReq.RemoteIP = connection.IPAddress(d.Get("remote_ip").(string))
	}
	if d.HasChange("remote_networks") {
		hasChange = true
		patchReq.RemoteNetworks = d.Get("remote_networks").(string)
	}
	if d.HasChange("local_networks") {
		hasChange = true
		patchReq.LocalNetworks = d.Get("local_networks").(string)
	}

	if hasChange {
		log.Printf("[INFO] Updating VPNSession with ID [%s]", d.Id())
		taskRef, err := service.PatchVPNSession(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating VPNSession with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for VPN session with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("psk") {
		updatePSKReq := ecloudservice.UpdateVPNSessionPreSharedKeyRequest{
			PSK: d.Get("psk").(string),
		}
		log.Printf("[DEBUG] Created UpdateVPNSessionPreSharedKeyRequest: %+v", updatePSKReq)

		taskRef, err := service.UpdateVPNSessionPreSharedKey(d.Id(), updatePSKReq)
		if err != nil {
			return fmt.Errorf("Error creating VPN session pre-shared key: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for VPN session with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceVPNSessionRead(d, meta)
}

func resourceVPNSessionDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing VPNSession with ID [%s]", d.Id())
	taskID, err := service.DeleteVPNSession(d.Id())
	if err != nil {
		return fmt.Errorf("Error VPNSession with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPNSession with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
