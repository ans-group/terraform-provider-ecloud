package ecloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					id := val.(string)
					fipAssignableResources := []string{"nic-"}

					prefixInSlice := func (slice []string, value string) bool {
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

func resourceFloatingIPCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateFloatingIPRequest{
		Name:  d.Get("name").(string),
		VPCID: d.Get("vpc_id").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
	}

	log.Printf("[DEBUG] Created CreateFloatingIPRequest: %+v", createReq)

	log.Print("[INFO] Creating Floating IP")
	taskRef, err := service.CreateFloatingIP(createReq)
	if err != nil {
		return fmt.Errorf("Error creating floating IP: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for floating IP with ID [%s] to be created: %s", d.Id(), err)
	}

	if r, ok := d.GetOk("resource_id"); ok {
		log.Printf("[DEBUG] Assigning floating IP with ID [%s] to resource [%s]", d.Id(), r.(string))

		assignFipReq := ecloudservice.AssignFloatingIPRequest{
			ResourceID: r.(string),
		}
		log.Printf("[DEBUG] Created AssignFloatingIPRequest: %+v", assignFipReq)

		taskID, err := service.AssignFloatingIP(d.Id(), assignFipReq)
		if err != nil {
			return fmt.Errorf("Error assigning floating IP: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
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
	d.Set("availability_zone_id", fip.AvailabilityZoneID)
	d.Set("resource_id", fip.ResourceID)

	return nil
}

func resourceFloatingIPUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating floating ip with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchFloatingIPRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchFloatingIP(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating floating ip with ID [%s]: %w", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	if d.HasChange("resource_id") {
		log.Printf("[INFO] Updating floating ip with ID [%s]", d.Id())

		oldVal, newVal := d.GetChange("resource_id")

		//if oldVal wasn't empty then floating ip needs unassigned first
		if oldVal.(string) != "" {
			log.Printf("[DEBUG] Unassigning floating IP with ID [%s]", d.Id())
			taskID, err := service.UnassignFloatingIP(d.Id())
			if err != nil {
				return fmt.Errorf("Error unassigning floating ip with ID [%s]: %w", d.Id(), err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(service, taskID),
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
			log.Printf("[DEBUG] Assigning floating IP with ID [%s] to resource [%s]", d.Id(), newVal.(string))

			assignFipReq := ecloudservice.AssignFloatingIPRequest{
				ResourceID: newVal.(string),
			}
			log.Printf("[DEBUG] Created AssignFloatingIPRequest: %+v", assignFipReq)

			taskID, err := service.AssignFloatingIP(d.Id(), assignFipReq)
			if err != nil {
				return fmt.Errorf("Error assigning floating IP: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Target:     []string{ecloudservice.TaskStatusComplete.String()},
				Refresh:    TaskStatusRefreshFunc(service, taskID),
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
	if _, ok := d.GetOk("resource_id"); ok {
		taskID, err := service.UnassignFloatingIP(d.Id())
		if err != nil {
			return fmt.Errorf("Error unassigning floating ip with ID [%s]: %w", d.Id(), err)
		}

		unassignStateConf := &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskID),
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
	taskID, err := service.DeleteFloatingIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.FloatingIPNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing floating ip with ID [%s]: %s", d.Id(), err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskID),
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
