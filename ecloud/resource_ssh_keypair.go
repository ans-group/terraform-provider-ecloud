package ecloud

import (
	"context"
	"fmt"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSshKeyPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyPairCreate,
		ReadContext:   resourceSshKeyPairRead,
		UpdateContext: resourceSshKeyPairUpdate,
		DeleteContext: resourceSshKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceSshKeyPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateSSHKeyPairRequest{
		PublicKey: d.Get("public_key").(string),
		Name:      d.Get("name").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateSSHKeyPairRequest: %+v", createReq))

	tflog.Info(ctx, "Creating SSH key pair")
	keyPairID, err := service.CreateSSHKeyPair(createReq)
	if err != nil {
		return diag.Errorf("Error creating ssh key pair: %s", err)
	}

	d.SetId(keyPairID)

	return resourceSshKeyPairRead(ctx, d, meta)
}

func resourceSshKeyPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Retrieving SSH key pair", map[string]interface{}{
		"id": d.Id(),
	})
	keyPair, err := service.GetSSHKeyPair(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.SSHKeyPairNotFoundError:
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", keyPair.Name)
	d.Set("public_key", keyPair.PublicKey)

	return nil
}

func resourceSshKeyPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	if d.HasChange("name") {
		tflog.Info(ctx, "Updating SSH key pair", map[string]interface{}{
			"id": d.Id(),
		})
		patchReq := ecloudservice.PatchSSHKeyPairRequest{
			Name: d.Get("name").(string),
		}

		err := service.PatchSSHKeyPair(d.Id(), patchReq)
		if err != nil {
			return diag.Errorf("Error updating ssh key pair with ID [%s]: %s", d.Id(), err)
		}
	}

	return resourceSshKeyPairRead(ctx, d, meta)
}

func resourceSshKeyPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Info(ctx, "Removing SSH key pair", map[string]interface{}{
		"id": d.Id(),
	})
	err := service.DeleteSSHKeyPair(d.Id())
	if err != nil {
		switch err.(type) {
		case *ecloudservice.SSHKeyPairNotFoundError:
			return nil
		default:
			return diag.Errorf("Error removing ssh key pair with ID [%s]: %s", d.Id(), err)
		}
	}

	return nil
}
