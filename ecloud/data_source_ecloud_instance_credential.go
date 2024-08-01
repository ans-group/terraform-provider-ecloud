package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstanceCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceCredentialsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"credential_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInstanceCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if username, ok := d.GetOk("username"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("username", connection.EQOperator, []string{username.(string)}))
	}
	if credentialID, ok := d.GetOk("credential_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{credentialID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	credentials, err := service.GetInstanceCredentials(d.Get("instance_id").(string), params)
	if err != nil {
		return diag.Errorf("Error retrieving instance credentials: %s", err)
	}

	if len(credentials) < 1 {
		return diag.Errorf("No credentials found with provided arguments")
	}

	if len(credentials) > 1 {
		return diag.Errorf("More than 1 credential found with provided arguments")
	}

	d.SetId(credentials[0].ID)
	d.Set("instance_id", credentials[0].ResourceID)
	d.Set("username", credentials[0].Username)
	d.Set("name", credentials[0].Name)
	d.Set("password", credentials[0].Password)

	return nil
}
