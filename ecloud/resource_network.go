package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateNetworkRequest{
		RouterID: d.Get("router_id").(string),
		Subnet:   d.Get("subnet").(string),
		Name:     d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateNetworkRequest: %+v", createReq)

	log.Print("[INFO] Creating network")
	networkID, err := service.CreateNetwork(createReq)
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	d.SetId(networkID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    NetworkSyncStatusRefreshFunc(service, networkID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for network with ID [%s] to return sync status of [%s]: %s", networkID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceNetworkRead(d, meta)
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving network with ID [%s]", d.Id())
	network, err := service.GetNetwork(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.NetworkNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("router_id", network.RouterID)
	d.Set("subnet", network.Subnet)
	d.Set("name", network.Name)

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchNetworkRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating network with ID [%s]", d.Id())
		err := service.PatchNetwork(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating network with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    NetworkSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 1 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for network with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceNetworkRead(d, meta)
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing network with ID [%s]", d.Id())
	err := service.DeleteNetwork(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing network with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    NetworkSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for network with ID [%s] to be deleted: %s", d.Id(), err)
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
