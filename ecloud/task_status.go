package ecloud

import (
	"fmt"
	"log"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TaskStatusRefreshFunc returns a function with StateRefreshFunc signature for use with StateChangeConf
func TaskStatusRefreshFunc(service ecloudservice.ECloudService, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		//check task status
		log.Printf("[DEBUG] Retrieving task status for taskID: [%s]", taskID)
		task, err := service.GetTask(taskID)
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] TaskID: %s has status: %s", task.ID, task.Status)

		if task.Status == ecloudservice.TaskStatusFailed {
			return nil, "", fmt.Errorf("Task with ID: %s has status of %s", task.ID, task.Status)
		}

		return "", task.Status.String(), nil
	}
}
