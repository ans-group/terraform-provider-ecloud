package ecloud

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInstanceScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceScriptCreate,
		ReadContext:   resourceInstanceScriptRead,
		DeleteContext: resourceInstanceScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"script": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceInstanceScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.ExecuteInstanceScriptRequest{
		Script:   d.Get("script").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	tflog.Debug(ctx, fmt.Sprintf("Created CreateInstanceScriptRequest: %+v", createReq))

	tflog.Info(ctx, "Executing InstanceScript")
	taskID, err := service.ExecuteInstanceScript(d.Get("instance_id").(string), createReq)
	if err != nil {
		return diag.Errorf("Error executing InstanceScript %s", err)
	}

	d.SetId(fmt.Sprintf("%d", rand.Int()))

	stateConf := &resource.StateChangeConf{
		Target:     []string{ecloudservice.TaskStatusComplete.String()},
		Refresh:    TaskStatusRefreshFunc(ctx, service, taskID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for InstanceScript with ID [%s] to run: %s", d.Id(), err)
	}

	return nil
}

func resourceInstanceScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
