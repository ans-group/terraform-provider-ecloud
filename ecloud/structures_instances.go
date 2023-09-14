package ecloud

import (
	"context"
	"fmt"

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
