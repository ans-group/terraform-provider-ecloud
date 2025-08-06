package ecloud

import (
	"context"
	"errors"
	"fmt"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext:   resourceTagRead,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateTagRequest{
		Name: d.Get("name").(string),
	}

	if scope, ok := d.GetOk("scope"); ok {
		createReq.Scope = scope.(string)
	}

	tflog.Debug(ctx, fmt.Sprintf("ecloud: creating tag with name [%s]", createReq.Name))

	tagID, err := service.CreateTag(createReq)
	if err != nil {
		return diag.Errorf("ecloud: error creating tag: %s", err)
	}

	d.SetId(tagID)

	return resourceTagRead(ctx, d, meta)
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Debug(ctx, fmt.Sprintf("ecloud: reading tag with ID [%s]", d.Id()))

	tag, err := service.GetTag(d.Id())
	if err != nil {
		var tagNotFoundError *ecloudservice.TagNotFoundError
		switch {
		case errors.As(err, &tagNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", tag.Name)
	d.Set("scope", tag.Scope)
	d.Set("created_at", tag.CreatedAt.String())
	d.Set("updated_at", tag.UpdatedAt.String())

	return nil
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	patchReq := ecloudservice.PatchTagRequest{}

	if d.HasChange("name") {
		patchReq.Name = d.Get("name").(string)
	}

	if d.HasChange("scope") {
		patchReq.Scope = d.Get("scope").(string)
	}

	tflog.Debug(ctx, fmt.Sprintf("ecloud: updating tag with ID [%s]", d.Id()))

	err := service.PatchTag(d.Id(), patchReq)
	if err != nil {
		return diag.Errorf("ecloud: error updating tag with ID [%s]: %s", d.Id(), err)
	}

	return resourceTagRead(ctx, d, meta)
}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	tflog.Debug(ctx, fmt.Sprintf("ecloud: deleting tag with ID [%s]", d.Id()))

	err := service.DeleteTag(d.Id())
	if err != nil {
		var tagNotFoundError *ecloudservice.TagNotFoundError
		switch {
		case errors.As(err, &tagNotFoundError):
			return nil
		default:
			return diag.Errorf("ecloud: error deleting tag with ID [%s]: %s", d.Id(), err)
		}
	}

	return nil
}
