package ecloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCCreate,
		Read:   resourceVPCRead,
		Update: resourceVPCUpdate,
		Delete: resourceVPCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
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

func resourceVPCCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVPCRequest{
		RegionID: d.Get("region_id").(string),
		Name:     d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateVPCRequest: %+v", createReq)

	log.Print("[INFO] Creating VPC")
	vpcID, err := service.CreateVPC(createReq)
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}

	d.SetId(vpcID)

	return resourceVPCRead(d, meta)
}

func resourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving VPC with ID [%s]", d.Id())
	vpc, err := service.GetVPC(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VPCNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("region_id", vpc.RegionID)
	d.Set("name", vpc.Name)

	return nil
}

func resourceVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		patchReq := ecloudservice.PatchVPCRequest{
			Name: d.Get("name").(string),
		}

		log.Printf("[INFO] Updating VPC with ID [%s]", d.Id())
		err := service.PatchVPC(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating VPC with ID [%s]: %w", d.Id(), err)
		}
	}

	return resourceVPCRead(d, meta)
}

func resourceVPCDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing VPC with ID [%s]", d.Id())
	err := service.DeleteVPC(d.Id())
	if err != nil {
		return fmt.Errorf("Error VPC with ID [%s]: %s", d.Id(), err)
	}

	return nil
}
