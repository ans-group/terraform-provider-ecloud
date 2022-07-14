package ecloud

import (
	"fmt"
	"log"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSshKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceSshKeyPairCreate,
		Read:   resourceSshKeyPairRead,
		Update: resourceSshKeyPairUpdate,
		Delete: resourceSshKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"public_key": {
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
	}
}

func resourceSshKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateSSHKeyPairRequest{
		PublicKey: d.Get("public_key").(string),
		Name:      d.Get("name").(string),
	}
	log.Printf("[DEBUG] Created CreateSSHKeyPairRequest: %+v", createReq)

	log.Print("[INFO] Creating ssh key pair")
	keyPairID, err := service.CreateSSHKeyPair(createReq)
	if err != nil {
		return fmt.Errorf("Error creating ssh key pair: %s", err)
	}

	d.SetId(keyPairID)

	return resourceSshKeyPairRead(d, meta)
}

func resourceSshKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Retrieving ssh key pair with ID [%s]", d.Id())
	keyPair, err := service.GetSSHKeyPair(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.SSHKeyPairNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", keyPair.Name)
	d.Set("public_key", keyPair.PublicKey)

	return nil
}

func resourceSshKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		log.Printf("[INFO] Updating ssh key pair with ID [%s]", d.Id())
		patchReq := ecloudservice.PatchSSHKeyPairRequest{
			Name: d.Get("name").(string),
		}

		err := service.PatchSSHKeyPair(d.Id(), patchReq)
		if err != nil {
			return fmt.Errorf("Error updating ssh key pair with ID [%s]: %w", d.Id(), err)
		}
	}

	return resourceSshKeyPairRead(d, meta)
}

func resourceSshKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	log.Printf("[INFO] Removing host group with ID [%s]", d.Id())
	err := service.DeleteSSHKeyPair(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.SSHKeyPairNotFoundError:
			return nil
		default:
			return fmt.Errorf("Error removing ssh key pair with ID [%s]: %s", d.Id(), err)
		}
	}

	return nil
}
