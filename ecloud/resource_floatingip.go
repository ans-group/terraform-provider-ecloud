package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceFloatingIPCreate,
		Read:   resourceFloatingIPRead,
		Update: resourceFloatingIPUpdate,
		Delete: resourceFloatingIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFloatingIPCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateFloatingIPRequest{
		Name:  d.Get("name").(string),
		VPCID: d.Get("vpc_id").(string),
	}

	log.Printf("[DEBUG] Created CreateFloatingIPRequest: %+v", createReq)

	log.Print("[INFO] Creating Floating IP")
	fipID, err := service.CreateFloatingIP(createReq)
	if err != nil {
		return fmt.Errorf("Error creating floating IP: %s", err)
	}

	d.SetId(fipID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    FloatingIPSyncStatusRefreshFunc(service, fipID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for floating IP with ID [%s] to be created: %s", d.Id(), err)
	}

	if r, ok := d.GetOk("instance_id"); ok {
		log.Printf("[DEBUG] Assigning floating IP with ID [%s]", d.Id())

		//retrieve instance nics
		nics, err := service.GetInstanceNICs(r.(string), connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Failed to retrieve instance nics: %w", err)
		}

		assignFipReq := ecloudservice.AssignFloatingIPRequest{
			ResourceID: nics[0].ID,
		}
		log.Printf("[DEBUG] Created AssignFloatingIPRequest: %+v", assignFipReq)

		err = service.AssignFloatingIP(d.Id(), assignFipReq)
		if err != nil {
			return fmt.Errorf("Error assigning floating IP: %s", err)
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for floating IP with ID [%s] to be assigned: %s", d.Id(), err)
		}
	}
	return resourceFloatingIPRead(d, meta)
}

func resourceFloatingIPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving floating ip with ID [%s]", d.Id())
	fip, err := service.GetFloatingIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FloatingIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", fip.VPCID)
	d.Set("name", fip.Name)
	d.Set("ip_address", fip.IPAddress)

	return nil
}

func resourceFloatingIPUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating floating ip with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchFloatingIPRequest{
			Name: d.Get("name").(string),
		}

		err := service.PatchFloatingIP(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating floating ip with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FloatingIPSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	if d.HasChange("instance_id") {
		log.Printf("[INFO] Updating floating ip with ID [%s]", d.Id())

		oldVal, newVal := d.GetChange("instance_id")

		//if oldVal wasn't empty then floating ip needs unassigned first
		if oldVal.(string) != "" {
			log.Printf("[DEBUG] Unassigning floating IP with ID [%s]", d.Id())
			err := service.UnassignFloatingIP(d.Id())
			if err != nil {
				return fmt.Errorf("Error unassigning floating ip with ID [%s]: %w", d.Id(), err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.SyncStatusComplete.String()},
				Refresh:    FloatingIPSyncStatusRefreshFunc(service, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      3 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %w", d.Id(), err)
			}
		}

		//Assign floating ip to new instance value if set
		if len(newVal.(string)) > 1 {
			log.Printf("[DEBUG] Assigning floating IP with ID [%s]", d.Id())

			//retrieve instance nics
			nics, err := service.GetInstanceNICs(newVal.(string), connection.APIRequestParameters{})
			if err != nil {
				return fmt.Errorf("Failed to retrieve instance nics: %w", err)
			}

			assignFipReq := ecloudservice.AssignFloatingIPRequest{
				ResourceID: nics[0].ID,
			}
			log.Printf("[DEBUG] Created AssignFloatingIPRequest: %+v", assignFipReq)

			err = service.AssignFloatingIP(d.Id(), assignFipReq)
			if err != nil {
				return fmt.Errorf("Error assigning floating IP: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.SyncStatusComplete.String()},
				Refresh:    FloatingIPSyncStatusRefreshFunc(service, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutCreate),
				Delay:      3 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf("Error waiting for floating IP with ID [%s] to be assigned: %s", d.Id(), err)
			}
		}
	}
	return resourceFloatingIPRead(d, meta)
}

func resourceFloatingIPDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	//first check if floating ip is assigned
	if _, ok := d.GetOk("instance_id"); ok {
		err := service.UnassignFloatingIP(d.Id())
		if err != nil {
			return fmt.Errorf("Error unassigning floating ip with ID [%s]: %w", d.Id(), err)
		}

		unassignStateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    FloatingIPSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = unassignStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %w", d.Id(), err)
		}
	}

	// Once unassigned - remove the floating ip resource
	log.Printf("[INFO] Removing floating ip with ID [%s]", d.Id())
	err := service.DeleteFloatingIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FloatingIPNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing floating ip with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    FloatingIPSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for floating ip with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}

func FloatingIPSyncStatusRefreshFunc(service ecloudservice.ECloudService, fipID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		fip, err := service.GetFloatingIP(fipID)
		if err != nil {
			if _, ok := err.(*ecloudservice.FloatingIPNotFoundError); ok {
				return fip, "Deleted", nil
			}
			return nil, "", err
		}

		if fip.Sync.Status == ecloudservice.SyncStatusFailed {
			return nil, "", fmt.Errorf("Failed to create/update floating ip - review logs")
		}

		return fip, fip.Sync.Status.String(), nil
	}
}
