package ecloud

import "log"

func expandCreateInstanceRequestImageData(rawData map[string]interface{}) map[string]interface{} {
	imageData := make(map[string]interface{})

	for k, v := range rawData {
		imageData[k] = v
	}

	log.Printf("[INFO] Image data: [%+v]", imageData)

	return imageData
}

func expandSshKeyPairIds(rawKeys []interface{}) []string {
	keyPairs := make([]string, len(rawKeys))

	for i, v := range rawKeys {
        keyPairs[i] = v.(string)
	}

	log.Printf("[INFO] SSH key pairs: [%+v]", keyPairs)

	return keyPairs
}
