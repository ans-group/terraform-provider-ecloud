package ecloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceLoadBalancerNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerNetworkCreate,
		Read:   resourceLoadBalancerNetworkRead,
		Update: resourceLoadBalancerNetworkUpdate,
		Delete: resourceLoadBalancerNetworkDelete,
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
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoadBalancerNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateLoadBalancerNetworkRequest{
		LoadBalancerID: d.Get("load_balancer_id").(string),
		NetworkID:      d.Get("network_id").(string),
		Name:           d.Get("name").(string),
	}

	log.Printf("[DEBUG] Created CreateLoadBalancerNetworkRequest: %+v", createReq)

	log.Print("[INFO] Creating LoadBalancerNetwork")
	taskRef, err := service.CreateLoadBalancerNetwork(createReq)
	if err != nil {
		return fmt.Errorf("Error creating loadbalancer network: %s", err)
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
		return fmt.Errorf("Error waiting for loadbalancer network with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceLoadBalancerNetworkRead(d, meta)
}

func resourceLoadBalancerNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving loadbalancer network with ID [%s]", d.Id())
	lbNetwork, err := service.GetLoadBalancerNetwork(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.LoadBalancerNetworkNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("network_id", lbNetwork.NetworkID)
	d.Set("name", lbNetwork.Name)
	d.Set("load_balancer_id", lbNetwork.LoadBalancerID)

	return nil
}

func resourceLoadBalancerNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating loadbalancer network with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchLoadBalancerNetworkRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchLoadBalancerNetwork(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating loadbalancer network with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for loadbalancer network with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceLoadBalancerNetworkRead(d, meta)
}

func resourceLoadBalancerNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing loadbalancer network with ID [%s]", d.Id())
	taskID, err := service.DeleteLoadBalancerNetwork(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.LoadBalancerNetworkNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing loadbalancer network with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for loadbalancer network with ID [%s] to be deleted: %s", d.Id(), err)
	}

	return nil
}
