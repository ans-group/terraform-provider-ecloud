package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostCreate,
		Read:   resourceHostRead,
		Update: resourceHostUpdate,
		Delete: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"host_group_id": {
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceHostCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateHostRequest{
		HostGroupID: d.Get("host_group_id").(string),
		Name:        d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateHostRequest: %+v", createReq)

	log.Print("[INFO] Creating Host Group")
	task, err := service.CreateHost(createReq)
	if err != nil {
		return fmt.Errorf("Error creating host group: %s", err)
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
		return fmt.Errorf("Error waiting for host group with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceHostRead(d, meta)
}

func resourceHostRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving host group with ID [%s]", d.Id())
	host, err := service.GetHost(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("host_group_id", host.HostGroupID)
	d.Set("name", host.Name)

	return nil
}

func resourceHostUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating host with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchHostRequest{
			Name: d.Get("name").(string),
		}

		task, err := service.PatchHost(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating host with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for host group with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceHostRead(d, meta)
}

func resourceHostDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing host group with ID [%s]", d.Id())
	taskID, err := service.DeleteHost(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.HostNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing host with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for host with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
