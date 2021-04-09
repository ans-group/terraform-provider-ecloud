package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouterCreate,
		Read:   resourceRouterRead,
		Update: resourceRouterUpdate,
		Delete: resourceRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			},
		},
	}
}

func resourceRouterCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateRouterRequest{
		VPCID: d.Get("vpc_id").(string),
		Name:  d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateRouterRequest: %+v", createReq)

	log.Print("[INFO] Creating Router")
	routerID, err := service.CreateRouter(createReq)
	if err != nil {
		return fmt.Errorf("Error creating router: %s", err)
	}

	d.SetId(routerID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    RouterSyncStatusRefreshFunc(service, routerID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for router with ID [%s] to return sync status of [%s]: %s", routerID, ecloudservice.SyncStatusComplete, err)
	}

	return resourceRouterRead(d, meta)
}

func resourceRouterRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving router with ID [%s]", d.Id())
	router, err := service.GetRouter(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.RouterNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", router.VPCID)
	d.Set("name", router.Name)

	return nil
}

func resourceRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchRouterRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating router with ID [%s]", d.Id())
		err := service.PatchRouter(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating router with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.SyncStatusComplete.String()},
			Refresh:    RouterSyncStatusRefreshFunc(service, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for router with ID [%s] to return sync status of [%s]: %s", d.Id(), ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceRouterRead(d, meta)
}

func resourceRouterDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing router with ID [%s]", d.Id())
	err := service.DeleteRouter(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing router with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    RouterSyncStatusRefreshFunc(service, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for router with ID [%s] to be deleted: %s", d.Id(), err)
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
