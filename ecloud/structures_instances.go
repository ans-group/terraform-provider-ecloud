package ecloud

import (
	"context"
	"fmt"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func expandCreateInstanceRequestImageData(ctx context.Context, rawData map[string]interface{}) map[string]interface{} {
	imageData := make(map[string]interface{})

	for k, v := range rawData {
		imageData[k] = v
	}

	tflog.Info(ctx, fmt.Sprintf("Image data: %+v", imageData))

	return imageData
}

func expandSshKeyPairIds(ctx context.Context, rawKeys []interface{}) []string {
	keyPairs := make([]string, len(rawKeys))

	for i, v := range rawKeys {
		keyPairs[i] = v.(string)
	}

	tflog.Info(ctx, fmt.Sprintf("SSH key pairs: %+v", keyPairs))

	return keyPairs
}

func flattenInstanceTags(tags []ecloudservice.ResourceTag) []interface{} {
	flattenedTags := make([]interface{}, len(tags))

	for i, tag := range tags {
		flattenedTag := make(map[string]interface{})
		flattenedTag["id"] = tag.ID
		flattenedTag["name"] = tag.Name
		flattenedTag["scope"] = tag.Scope
		flattenedTags[i] = flattenedTag
	}

	return flattenedTags
}
