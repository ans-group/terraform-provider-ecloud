package ecloud

import (
	"fmt"
	"log"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	instanceID := d.Get("instance_id").(string)
	createReq := ecloudservice.CreateInstanceImageRequest{
		Name: d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateImageRequest: %+v", createReq)

	log.Print("[INFO] Creating image")
	taskRef, err := service.CreateInstanceImage(instanceID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating image: %s", err)
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
		return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.TaskStatusComplete, err)
	}

	return resourceImageRead(d, meta)
}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving image with ID [%s]", d.Id())
	image, err := service.GetImage(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.ImageNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", image.VPCID)
	d.Set("name", image.Name)
	d.Set("availability_zone_id", image.AvailabilityZoneID)

	return nil
}

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	hasChange := false
	patchReq := ecloudservice.UpdateImageRequest{}

	if d.HasChange("name") {
		hasChange = true
		patchReq.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[INFO] Updating image with ID [%s]", d.Id())
		taskRef, err := service.UpdateImage(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating image with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for task with ID [%s] to return task status of [%s]: %s", taskRef.TaskID, ecloudservice.SyncStatusComplete, err)
		}
	}

	return resourceImageRead(d, meta)
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing image with ID [%s]", d.Id())
	taskID, err := service.DeleteImage(d.Id())
	if err != nil {
		return fmt.Errorf("Error removing image with ID [%s]: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.SyncStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for image with ID [%s] to be deleted: %s", taskID, err)
	}

	return nil
}
