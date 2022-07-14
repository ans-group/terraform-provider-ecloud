package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSshKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSshKeyPairRead,

		Schema: map[string]*schema.Schema{
			"ssh_keypair_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSshKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("ssh_keypair_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	keyPairs, err := service.GetSSHKeyPairs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving ssh keypair: %s", err)
	}

	if len(keyPairs) < 1 {
		return errors.New("No ssh keypairs found with provided arguments")
	}

	if len(keyPairs) > 1 {
		return errors.New("More than 1 host found with provided arguments")
	}

	d.SetId(keyPairs[0].ID)
	d.Set("name", keyPairs[0].Name)
	d.Set("public_key", keyPairs[0].PublicKey)

	return nil
}
