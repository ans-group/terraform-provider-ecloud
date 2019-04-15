package ecloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVirtualMachineTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineTagCreate,
		Read:   resourceVirtualMachineTagRead,
		Update: resourceVirtualMachineTagUpdate,
		Delete: resourceVirtualMachineTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"virtualmachine_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVirtualMachineTagCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmid := d.Get("virtualmachine_id").(int)
	tagKey := d.Get("key").(string)
	tagValue := d.Get("value").(string)

	createReq := ecloudservice.CreateTagRequest{
		Key:   tagKey,
		Value: tagValue,
	}
	log.Printf("Created CreateTagRequest: %+v", createReq)

	log.Print("Creating virtual machine tag")
	err := service.CreateVirtualMachineTag(vmid, createReq)
	if err != nil {
		return fmt.Errorf("Error creating virtual machine tag: %s", err)
	}

	d.SetId(tagKey)

	return resourceVirtualMachineTagRead(d, meta)
}

func resourceVirtualMachineTagRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmid := d.Get("virtualmachine_id").(int)
	tagKey := d.Id()

	log.Printf("Retrieving tag with key [%s] for virtual machine with ID [%d]", tagKey, vmid)
	tag, err := service.GetVirtualMachineTag(vmid, tagKey)
	if err != nil {
		switch err.(type) {
		case *ecloudservice.TagNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("key", tag.Key)
	d.Set("value", tag.Value)

	return nil
}

func resourceVirtualMachineTagUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmid := d.Get("virtualmachine_id").(int)
	tagKey := d.Id()

	if d.HasChange("value") {
		patchReq := ecloudservice.PatchTagRequest{
			Value: d.Get("value").(string),
		}

		log.Printf("Updating tag with key [%s] for virtual machine with ID [%d]", tagKey, vmid)
		err := service.PatchVirtualMachineTag(vmid, tagKey, patchReq)
		if err != nil {
			return fmt.Errorf("Error updating tag with key [%s] for virtual machine with ID [%d]", tagKey, vmid)
		}
	}

	return resourceVirtualMachineTagRead(d, meta)
}

func resourceVirtualMachineTagDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmid := d.Get("virtualmachine_id").(int)
	tagKey := d.Id()

	log.Printf("Removing tag with key [%s] for virtual machine with ID [%d]", tagKey, vmid)
	err := service.DeleteVirtualMachineTag(vmid, tagKey)
	if err != nil {
		return fmt.Errorf("Error removing tag with key [%s] for virtual machine with ID [%d]: %s", tagKey, vmid, err)
	}

	return nil
}
