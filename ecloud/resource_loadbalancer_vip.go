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

func resourceLoadBalancerVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerVipCreate,
		Read:   resourceLoadBalancerVipRead,
		Update: resourceLoadBalancerVipUpdate,
		Delete: resourceLoadBalancerVipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allocate_floating_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLoadBalancerVipCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVIPRequest{
		LoadBalancerID: d.Get("load_balancer_id").(string),
		AllocateFloatingIP:  d.Get("allocate_floating_ip").(bool),
		Name:           d.Get("name").(string),
	}

	log.Printf("[DEBUG] Created CreateVIPRequest: %+v", createReq)

	log.Print("[INFO] Creating LoadBalancer VIP")
	taskRef, err := service.CreateVIP(createReq)
	if err != nil {
		return fmt.Errorf("Error creating loadbalancer vip: %s", err)
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
		return fmt.Errorf("Error waiting for loadbalancer vip with ID [%s] to be created: %s", d.Id(), err)
	}

	return resourceLoadBalancerVipRead(d, meta)
}

func resourceLoadBalancerVipRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving loadbalancer vip with ID [%s]", d.Id())
	lbVip, err := service.GetVIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VIPNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", lbVip.Name)
	d.Set("load_balancer_id", lbVip.LoadBalancerID)

	
	if d.Get("floating_ip_id").(string) == "" && d.Get("allocate_floating_ip").(bool) {
		//we need to use the IP ID from the vip
		//and filter for the floating IP ID
		params := connection.APIRequestParameters{}
		params.WithFilter(*connection.NewAPIRequestFiltering("resource_id", connection.EQOperator, []string{lbVip.IPAddressID}))

		fips, err := service.GetFloatingIPs(params)
		if err != nil {
			return fmt.Errorf("Failed to retrieve floating IPs: %w", err)
		}

		if len(fips) != 1 {
			return fmt.Errorf("Unexpected number of floating IPs allocated to VIP")
		}

		d.Set("floating_ip_id", fips[0].ID)
	}

	return nil
}

func resourceLoadBalancerVipUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating loadbalancer vip with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchVIPRequest{
			Name: d.Get("name").(string),
		}

		taskRef, err := service.PatchVIP(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating loadbalancer vip with ID [%s]: %w", d.Id(), err)
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
			return fmt.Errorf("Error waiting for loadbalancer vip with ID [%s] to return task status of [%s]: %s", d.Id(), ecloudservice.TaskStatusComplete, err)
		}
	}

	return resourceLoadBalancerVipRead(d, meta)
}

func resourceLoadBalancerVipDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing loadbalancer network with ID [%s]", d.Id())
	taskID, err := service.DeleteVIP(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VIPNotFoundError:
			log.Printf("[DEBUG] loadbalancer VIP with ID [%s] not found. Continuing.", d.Id())
		default:
			return fmt.Errorf("Error removing loadbalancer vip with ID [%s]: %s", d.Id(), err)
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
		return fmt.Errorf("Error waiting for loadbalancer vip with ID [%s] to be deleted: %s", d.Id(), err)
	}

	//remove floating ip if set
	if len(d.Get("floating_ip_id").(string)) > 1 {
		fip := d.Get("floating_ip_id").(string)

		log.Printf("[DEBUG] Unassigning floating ip with ID [%s]", fip)

		taskID, err := service.UnassignFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				log.Printf("[DEBUG] Floating IP with ID [%s] not found. Skipping unassign.", fip)
			default:
				return fmt.Errorf("Error unassigning floating ip with ID [%s]: %s", fip, err)
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
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to be unassigned: %w", d.Id(), err)
		}

		log.Printf("[DEBUG] Removing floating ip with ID [%s]", fip)

		taskID, err = service.DeleteFloatingIP(fip)
		if err != nil {
			switch err.(type) {
			case *ecloudservice.FloatingIPNotFoundError:
				log.Printf("[DEBUG] Floating IP with ID [%s] not found. Skipping delete.", fip)
			default:
				return fmt.Errorf("Error removing floating ip with ID [%s]: %s", fip, err)
			}
		}

		stateConf = &resource.StateChangeConf{
			Target:     []string{ecloudservice.TaskStatusComplete.String()},
			Refresh:    TaskStatusRefreshFunc(service, taskID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for floating ip with ID [%s] to be removed: %w", d.Id(), err)
		}
	}

	return nil
}
