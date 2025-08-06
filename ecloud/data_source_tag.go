package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTag() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagRead,

		Schema: map[string]*schema.Schema{
			"tag_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("tag_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if scope, ok := d.GetOk("scope"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("scope", connection.EQOperator, []string{scope.(string)}))
	}

	tags, err := service.GetTags(params)
	if err != nil {
		return diag.Errorf("ecloud: error retrieving tags: %s", err)
	}

	if len(tags) < 1 {
		return diag.Errorf("ecloud: no tags found with provided arguments")
	}

	if len(tags) > 1 {
		return diag.Errorf("ecloud: more than 1 tag found with provided arguments")
	}

	d.SetId(tags[0].ID)
	d.Set("name", tags[0].Name)
	d.Set("scope", tags[0].Scope)
	d.Set("created_at", tags[0].CreatedAt.String())
	d.Set("updated_at", tags[0].UpdatedAt.String())

	return nil
}
