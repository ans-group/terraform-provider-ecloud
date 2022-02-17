package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,
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
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_spec_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "network_ids": {
			// 	Type:     schema.TypeSet,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// 	Optional: true,
			// 	Set:      schema.HashString,
			// },
		},
	}
}

func resourceLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateLoadBalancerRequest{
		VPCID:              d.Get("vpc_id").(string),
		AvailabilityZoneID: d.Get("availability_zone_id").(string),
		LoadBalancerSpecID: d.Get("load_balancer_spec_id").(string),
	//	NetworkIDs:         expandNetworkIds(d.Get("network_ids").(*schema.Set).List()),
		Name:               d.Get("name").(string),
	}

	log.Printf("[DEBUG] Created CreateLoadBalancerRequest: %+v", createReq)

	log.Print("[INFO] Creating LoadBalancer")
	taskRef, err := service.CreateLoadBalancer(createReq)
	if err != nil {
		return fmt.Errorf("Error creating loadbalancer: %s", err)
	}

	d.SetId(taskRef.ResourceID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(service, taskRef.TaskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for loadbalancer with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving loadbalancer with ID [%s]", d.Id())
	lb, err := service.GetLoadBalancer(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.LoadBalancerNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("vpc_id", lb.VPCID)
	d.Set("name", lb.Name)
	d.Set("config_id", lb.ConfigID)
	d.Set("availability_zone_id", lb.AvailabilityZoneID)
	d.Set("load_balancer_spec_id", lb.LoadBalancerSpecID)
	

	return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating loadbalancer with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchLoadBalancerRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchLoadBalancer(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating loadbalancer with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for loadbalancer with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing loadbalancer with ID [%s]", d.Id())
	taskID, err := service.DeleteLoadBalancer(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.LoadBalancerNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing loadbalancer with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for loadbalancer with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
