package ecloud

import (
	"context"
	"fmt"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TaskStatusRefreshFunc returns a function with StateRefreshFunc signature for use with StateChangeConf
func TaskStatusRefreshFunc(ctx context.Context, service ecloudservice.ECloudService, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tflog.Debug(ctx, "Retrieving task status", map[string]interface{}{
			"task_id": taskID,
		})
		task, err := service.GetTask(taskID)
		if err != nil {
			return nil, "", err
		}

		tflog.Debug(ctx, "Retrieved task status", map[string]interface{}{
			"task_id":     task.ID,
			"task_status": task.Status,
		})

		if task.Status == ecloudservice.TaskStatusFailed {
			return nil, "", fmt.Errorf("Task with ID: %s has status of %s", task.ID, task.Status)
		}

		return "", task.Status.String(), nil
	}
}
