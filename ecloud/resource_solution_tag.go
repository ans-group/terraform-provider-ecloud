package ecloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceSolutionTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolutionTagCreate,
		Read:   resourceSolutionTagRead,
		Update: resourceSolutionTagUpdate,
		Delete: resourceSolutionTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"solution_id": &schema.Schema{
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

func resourceSolutionTagCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	solutionID := d.Get("solution_id").(int)
	tagKey := d.Get("key").(string)
	tagValue := d.Get("value").(string)

	createReq := ecloudservice.CreateTagRequest{
		Key:   tagKey,
		Value: tagValue,
	}
	log.Printf("Created CreateTagRequest: %+v", createReq)

	log.Print("Creating solution tag")
	err := service.CreateSolutionTag(solutionID, createReq)
	if err != nil {
		return fmt.Errorf("Error creating solution tag: %s", err)
	}

	d.SetId(tagKey)

	return resourceSolutionTagRead(d, meta)
}

func resourceSolutionTagRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	solutionID := d.Get("solution_id").(int)
	tagKey := d.Id()

	log.Printf("Retrieving tag with key [%s] for solution with ID [%d]", tagKey, solutionID)
	tag, err := service.GetSolutionTag(solutionID, tagKey)
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

func resourceSolutionTagUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	solutionID := d.Get("solution_id").(int)
	tagKey := d.Id()

	if d.HasChange("value") {
		patchReq := ecloudservice.PatchTagRequest{
			Value: d.Get("value").(string),
		}

		log.Printf("Updating tag with key [%s] for solution with ID [%d]", tagKey, solutionID)
		err := service.PatchSolutionTag(solutionID, tagKey, patchReq)
		if err != nil {
			return fmt.Errorf("Error updating tag with key [%s] for solution with ID [%d]", tagKey, solutionID)
		}
	}

	return resourceSolutionTagRead(d, meta)
}

func resourceSolutionTagDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	solutionID := d.Get("solution_id").(int)
	tagKey := d.Id()

	log.Printf("Removing tag with key [%s] for solution with ID [%d]", tagKey, solutionID)
	err := service.DeleteSolutionTag(solutionID, tagKey)
	if err != nil {
		return fmt.Errorf("Error removing tag with key [%s] for solution with ID [%d]: %s", tagKey, solutionID, err)
	}

	return nil
}
